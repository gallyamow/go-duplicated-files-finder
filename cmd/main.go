package main

import (
	"fmt"
	"github.com/gallyamow/go-duplicated-files-finder/internal/config"
	"github.com/gallyamow/go-duplicated-files-finder/internal/scanner"
	"log"
	"os"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Config: %+v\n", cfg)

	files, err := scanner.ScanDir(cfg.Path, cfg.MinSize)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d files\n", len(files))
}
