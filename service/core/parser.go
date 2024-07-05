package core

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

type PagePredicate = func(int) bool

func ParseManifest(manifestUrl, outFile string, prepareationChan <-chan error, pages PagePredicate) error {
	resp, err := http.Get(manifestUrl)
	if err != nil {
		return fmt.Errorf("error getting manifest: %w", err)
	}
	defer resp.Body.Close()
	manifest, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading manifest: %w", err)
	}

	imageURLs, err := parseImageURLs(manifest, pages)
	if err != nil {
		return fmt.Errorf("error parsing image URLs: %w", err)
	}
	fmt.Println("Extracted images:", len(imageURLs))

	dataDir := filepath.Join("data", strings.TrimSuffix(filepath.Base(outFile), filepath.Ext(outFile)))
	defer os.RemoveAll(dataDir)
	os.MkdirAll(dataDir, 0755)

	files, err := downloadImages(dataDir, imageURLs)
	if err != nil {
		return fmt.Errorf("error downloading images: %w", err)
	}

	if prepareationChan != nil {
		err = <-prepareationChan
		if err != nil {
			return err
		}
	}

	files = lo.Filter(files, func(f string, _ int) bool {
		return len(f) > 0
	})
	if len(files) == 0 {
		return fmt.Errorf("no images downloaded")
	}
	err = pdfcpu.ImportImagesFile(files, outFile, nil, nil)
	if err != nil {
		return fmt.Errorf("error creating PDF: %w", err)
	}

	fmt.Println("Done! ", outFile)
	return nil
}

func ParsePenn(catalogID string, pages PagePredicate) (string, error) {
	// Example: https://colenda.library.upenn.edu/catalog/81431-p3hk28

	fmt.Println("Catalog ID:", catalogID)
	outFile := fmt.Sprintf("%s.pdf", catalogID)

	headerCh := writeHeader(outFile, func() (Header, error) {
		return createPennHeader(catalogID)
	})

	return outFile, ParseManifest(fmt.Sprintf(pennManifestUrlFmt, catalogID), outFile, headerCh, pages)
}

func ParseShakespeare(catalogID string, pages PagePredicate) (string, error) {
	// Example: https://digitalcollections.folger.edu/bib244741-309974-lb41

	fmt.Println("Catalog ID:", catalogID)
	outFile := fmt.Sprintf("%s.pdf", catalogID)

	header, err := createShakespearHeader(catalogID)
	if err != nil {
		return "", fmt.Errorf("error creating header: %w", err)
	}

	headerCh := writeHeader(outFile, func() (Header, error) {
		return header, nil
	})

	return outFile, ParseManifest(header.ManifestURL, outFile, headerCh, pages)
}

func parseImageURLs(content []byte, pages PagePredicate) ([]string, error) {
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
		if pages != nil && !pages(i) {
			continue
		}
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

func writeHeader(outFilePath string, action func() (Header, error)) <-chan error {
	ch := make(chan error, 1)
	go func() {
		header, err := action()
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
	if len(header.Links) > 0 {
		b.WriteString("\nLinks:\n")
		for _, l := range header.Links {
			b.WriteString(l)
			b.WriteString("\n")
		}
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
