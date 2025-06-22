package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/boss-net/api/boss-plugin/internal/server/controllers/generator"
)

func main() {
	// Parse command line flags
	flag.Parse()

	// Generate all files
	if err := generator.GenerateAll(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
