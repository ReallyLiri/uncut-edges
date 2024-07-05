package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	concurrency = 10
	imageExt    = "jpg"
)

func downloadImages(downloadDir string, imageURLs []string) ([]string, error) {
	var wg sync.WaitGroup
	wg.Add(len(imageURLs))
	sem := make(chan struct{}, concurrency)
	defer close(sem)
	outFiles := make([]string, len(imageURLs))
	for i, url := range imageURLs {
		filePath := filepath.Join(downloadDir, fmt.Sprintf("img_%d.%s", i+1, imageExt))
		go func() {
			if success := downloadImage(url, filePath, &wg, sem); success {
				outFiles[i] = filePath
			}
		}()
	}
	wg.Wait()
	return outFiles, nil
}

func downloadImage(url string, filePath string, wg *sync.WaitGroup, sem chan struct{}) bool {
	defer func() {
		<-sem
		wg.Done()
	}()
	sem <- struct{}{}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading '%s': %v\n", url, err)
		return false
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error downloading '%s' with status %s:\n%s\n", url, resp.Status, body)
		return false
	}
	defer resp.Body.Close()
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filePath, err)
		return false
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", filePath, err)
		return false
	}

	fmt.Printf("Downloaded: %s\n", filePath)
	return true
}
