package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

func parseShakespeareHandler(w http.ResponseWriter, r *http.Request) {
	catalogID := r.URL.Path[len(shakespeareEndpoint):]
	if catalogID == "" {
		http.Error(w, "Missing play parameter", http.StatusBadRequest)
		return
	}

	filePath, err := core.ParseShakespeare(catalogID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing play: %v", err), http.StatusInternalServerError)
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

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(pdfFilePath)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error streaming file: %v", err), http.StatusInternalServerError)
		return
	}
}
