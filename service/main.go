package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <CatalogID>")
		return
	}
	catalogID := os.Args[1]

	// Example: https://colenda.library.upenn.edu/catalog/81431-p3hk28
	_, err := parse(catalogID)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
