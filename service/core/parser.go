package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

func Parse(catalogID string) (string, error) {
	fmt.Println("Catalog ID:", catalogID)

	resp, err := http.Get("https://colenda.library.upenn.edu/phalt/iiif/2/" + catalogID + "/manifest")
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

	outFile := fmt.Sprintf("%s.pdf", catalogID)
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
