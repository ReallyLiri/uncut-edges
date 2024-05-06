package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/reallyliri/penn-scans-parser/core"
)

const (
	parseEndpoint = "/parse/"
)

func parseCatalogHandler(w http.ResponseWriter, r *http.Request) {
	catalogID := r.URL.Path[len(parseEndpoint):]
	if catalogID == "" {
		http.Error(w, "Missing catalogID parameter", http.StatusBadRequest)
		return
	}

	filePath, err := core.Parse(catalogID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing catalog: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(filePath)

	f, err := os.Open(filePath)
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
