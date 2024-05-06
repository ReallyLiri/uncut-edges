package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

func Parse(catalogID string) (string, error) {
	fmt.Println("Catalog ID:", catalogID)
	outFile := fmt.Sprintf("%s.pdf", catalogID)

	headerCh := writeHeader(catalogID, outFile)

	resp, err := http.Get(fmt.Sprintf(manifestUrlFmt, catalogID))
	if err != nil {
		return "", fmt.Errorf("error getting manifest: %w", err)
	}
	defer resp.Body.Close()
	manifest, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading manifest: %w", err)
	}

	imageURLs, err := parseImageURLs(manifest)
	if err != nil {
		return "", fmt.Errorf("error parsing image URLs: %w", err)
	}
	fmt.Println("Extracted images:", len(imageURLs))

	dataDir := filepath.Join("data", catalogID)
	defer os.RemoveAll(dataDir)
	os.MkdirAll(dataDir, 0755)

	files, err := downloadImages(dataDir, imageURLs)
	if err != nil {
		return "", fmt.Errorf("error downloading images: %w", err)
	}

	err = <-headerCh
	if err != nil {
		return "", err
	}

	pdfcpu.ImportImagesFile(files, outFile, nil, nil)

	fmt.Println("Done! ", outFile)
	return outFile, nil
}

func parseImageURLs(content []byte) ([]string, error) {
	var data Manifest
	err := json.Unmarshal(content, &data)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}
	if len(data.Sequences) != 1 {
		return nil, fmt.Errorf("there should be exactly one sequence in the JSON file")
	}
	firstSequence := data.Sequences[0]
	canvases := firstSequence.Canvases
	urls := []string{}
	for i, canvas := range canvases {
		if len(canvas.Images) == 0 {
			return nil, fmt.Errorf("there are no images in the canvas of index %d", i)
		}
		firstImage := canvas.Images[0]
		fmt.Printf("Resource ID: %s\n", firstImage.Resource.ID)
		u := strings.Replace(firstImage.Resource.ID, "/full/!200,200/0/default.jpg", "/full/2000,/0/default.jpg", 1)
		urls = append(urls, u)
	}
	return urls, nil
}

func writeHeader(catalogID, outFilePath string) <-chan error {
	ch := make(chan error, 1)
	go func() {
		header, err := createHeaderFile(catalogID)
		if err != nil {
			ch <- err
			return
		}

		f, err := os.Create(outFilePath)
		if err != nil {
			ch <- fmt.Errorf("error creating file: %w", err)
			return
		}
		defer f.Close()

		ch <- formatHeaderJson(f, header)
	}()
	return ch
}

func formatHeaderJson(w io.Writer, header Header) error {
	b := strings.Builder{}
	b.WriteString("\n")
	b.WriteString(header.Title)
	b.WriteString("\n")
	b.WriteString("Catalog ID: ")
	b.WriteString(header.CatalogID)
	b.WriteString("\n\n")
	for _, p := range header.Properties {
		b.WriteString(p.Key)
		b.WriteString(" ")
		b.WriteString(p.Value)
		b.WriteString("\n")
	}
	b.WriteString("\nLinks:\n")
	for _, l := range header.Links {
		b.WriteString(l)
		b.WriteString("\n")
	}
	s := b.String()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.MoveTo(0, 10)
	pdf.SetFont("Arial", "", 12)
	pdf.SetLeftMargin(10)
	pdf.SetRightMargin(10)
	pdf.SetAutoPageBreak(true, 10)
	pdf.MultiCell(0.0, 6, s, "", "L", false)
	err := pdf.Output(w)
	if err != nil {
		return fmt.Errorf("error writing PDF: %w", err)
	}
	return nil
}
