package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/reallyliri/penn-scans-parser/core"
)

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

const (
	usage = `Usage: go run main.go <serve/parse> <ManifestURL>`
)

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
		manifestURL := os.Args[2]
		outFile := "out.pdf"
		err := core.ParseManifest(manifestURL, outFile, nil)
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
	http.HandleFunc(parseEndpoint, parseManifestHandler)
	http.HandleFunc(parsePennEndpoint, parsePennHandler)
	http.HandleFunc(shakespeareEndpoint, parseShakespeareHandler)

	fmt.Printf("Server is listening on port %s...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
