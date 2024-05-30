package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/reallyliri/penn-scans-parser/core"
)

const (
	parseEndpoint       = "/parse/"
	parsePennEndpoint   = parseEndpoint + "penn/"
	shakespeareEndpoint = parseEndpoint + "shakespeare/"
	vatlibEndpoint      = parseEndpoint + "vatlib/"

	allowedDomain = "https://uncut-edges.netlify.app" // "http://localhost:5173"

	pagesQueryParam = "pages"
)

func parseManifestHandler(w http.ResponseWriter, r *http.Request) {
	manifestUrl := r.URL.Path[len(parseEndpoint):]
	if manifestUrl == "" {
		http.Error(w, "Missing manifestUrl parameter", http.StatusBadRequest)
		return
	}

	pages, err := getPages(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing pages: %v", err), http.StatusBadRequest)
		return
	}

	outFile := "out.pdf"
	err = core.ParseManifest(manifestUrl, outFile, nil, pages)
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

	pages, err := getPages(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing pages: %v", err), http.StatusBadRequest)
		return
	}

	filePath, err := core.ParsePenn(catalogID, pages)
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

	pages, err := getPages(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing pages: %v", err), http.StatusBadRequest)
		return
	}

	filePath, err := core.ParseShakespeare(catalogID, pages)
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
	w.Header().Set("Access-Control-Allow-Origin", allowedDomain)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Expose-Headers", "*")

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error streaming file: %v", err), http.StatusInternalServerError)
		return
	}
}

func getPages(r *http.Request) (core.PagePredicate, error) {
	pages := r.URL.Query().Get(pagesQueryParam)
	if pages == "" {
		return nil, nil
	}
	ranges, err := parsePageRanges(pages)
	if err != nil {
		return nil, fmt.Errorf("error parsing pages '%s': %w", pages, err)
	}
	return func(i int) bool {
		for _, r := range ranges {
			if i >= int(r.Start) && i <= int(r.End) {
				return true
			}
		}
		return false
	}, nil
}

func parsePageRanges(pageRanges string) ([]core.Range, error) {
	var ranges []core.Range
	parts := strings.Split(pageRanges, ",")

	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("start of range is greater than end: %s", part)
			}

			ranges = append(ranges, core.Range{Start: start, End: end})
		} else {
			page, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", part)
			}
			ranges = append(ranges, core.Range{Start: page, End: page})
		}
	}

	return ranges, nil
}
