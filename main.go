package main

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <CatalogID>")
		return
	}
	catalogID := os.Args[1]

	// Example: https://colenda.library.upenn.edu/catalog/81431-p3hk28
	fmt.Println("Catalog ID:", catalogID)

	resp, err := http.Get("https://colenda.library.upenn.edu/phalt/iiif/2/" + catalogID + "/manifest")
	if err != nil {
		fmt.Println("Error getting manifest:", err)
		return
	}
	defer resp.Body.Close()
	manifest, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading manifest:", err)
		return
	}

	imageURLs, err := parseImageURLs(manifest)
	if err != nil {
		fmt.Println("Error parsing image URLs:", err)
		return
	}
	fmt.Println("Extracted images:", len(imageURLs))

	dataDir := filepath.Join("data", catalogID)
	defer os.RemoveAll(dataDir)
	os.MkdirAll(dataDir, 0755)

	files, err := downloadImages(dataDir, imageURLs)
	if err != nil {
		fmt.Println("Error downloading images:", err)
		return
	}

	outFile := fmt.Sprintf("%s.pdf", catalogID)
	pdfcpu.ImportImagesFile(files, outFile, nil, nil)

	fmt.Println("Done! ", outFile)
}

func parseImageURLs(content []byte) ([]string, error) {
	var data JSONData
	err := json.Unmarshal(content, &data)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}
	if len(data.Sequences) != 1 {
		return nil, fmt.Errorf("there should be exactly one sequence in the JSON file")
	}
	firstSequence := data.Sequences[0]
	canvases := firstSequence.Canvases
	jpegs := []string{}
	for i, canvas := range canvases {
		if len(canvas.Images) == 0 {
			return nil, fmt.Errorf("there are no images in the canvas of index %d", i)
		}
		firstImage := canvas.Images[0]
		fmt.Printf("Resource ID: %s\n", firstImage.Resource.ID)
		u := strings.Replace(firstImage.Resource.ID, "/full/!200,200/0/default.jpg", "/full/2000,/0/default.jpg", 1)
		jpegs = append(jpegs, u)
	}
	return jpegs, nil
}
