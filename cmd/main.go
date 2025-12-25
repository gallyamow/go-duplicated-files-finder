package main

import (
	"context"
	"fmt"
	"github.com/gallyamow/go-duplicated-files-finder/internal/config"
	"github.com/gallyamow/go-duplicated-files-finder/internal/hasher"
	"github.com/gallyamow/go-duplicated-files-finder/internal/scanner"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	fmt.Println(cfg)

	files, err := scanner.ScanDir(cfg.Path, cfg.MinSize)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d files\n", len(files))

	hashes := hasher.HashFiles(ctx, files, cfg.Algo, cfg.Workers)
	fmt.Printf("Hashed %d files\n", len(hashes))
}
