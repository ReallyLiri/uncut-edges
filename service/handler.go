package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/reallyliri/penn-scans-parser/core"
)

const (
	parseEndpoint       = "/parse/"
	parsePennEndpoint   = parseEndpoint + "penn/"
	shakespeareEndpoint = parseEndpoint + "shakespeare/"
	vatlibEndpoint      = parseEndpoint + "vatlib/"
)

func parseManifestHandler(w http.ResponseWriter, r *http.Request) {
	manifestUrl := r.URL.Path[len(parseEndpoint):]
	if manifestUrl == "" {
		http.Error(w, "Missing manifestUrl parameter", http.StatusBadRequest)
		return
	}

	outFile := "out.pdf"
	err := core.ParseManifest(manifestUrl, outFile, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing manifest: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(outFile)

	writeResponse(w, outFile)
}

func parsePennHandler(w http.ResponseWriter, r *http.Request) {
	catalogID := r.URL.Path[len(parsePennEndpoint):]
	if catalogID == "" {
		http.Error(w, "Missing catalogID parameter", http.StatusBadRequest)
		return
	}

	filePath, err := core.ParsePenn(catalogID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing catalog: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(filePath)

	writeResponse(w, filePath)
}

func writeResponse(w http.ResponseWriter, pdfFilePath string) {
	f, err := os.Open(pdfFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error streaming file: %v", err), http.StatusInternalServerError)
		return
	}
}
