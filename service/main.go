package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/reallyliri/penn-scans-parser/core"
)

const (
	port  = 8080
	usage = `Usage: go run main.go <serve/parse> <CatalogID>`
)

// Example: https://colenda.library.upenn.edu/catalog/81431-p3hk28

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	action := os.Args[1]
	switch action {
	case "serve":
		serve()
	case "parse":
		if len(os.Args) < 3 {
			fmt.Println(usage)
			os.Exit(1)
		}
		catalogID := os.Args[2]
		_, err := core.Parse(catalogID)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Println(usage)
		os.Exit(1)
	}
}

func serve() {
	http.HandleFunc(parseEndpoint, parseCatalogHandler)
	fmt.Printf("Server is listening on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
