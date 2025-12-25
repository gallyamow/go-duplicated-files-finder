package main

import (
	"context"
	"fmt"
	"github.com/gallyamow/go-duplicated-files-finder/internal/config"
	"github.com/gallyamow/go-duplicated-files-finder/internal/finder"
	"github.com/gallyamow/go-duplicated-files-finder/internal/hasher"
	"github.com/gallyamow/go-duplicated-files-finder/internal/printer"
	"github.com/gallyamow/go-duplicated-files-finder/internal/scanner"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var version = "unknown"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed parse config:", err)
		os.Exit(1)
	}

	if err := hasher.ValidateAlgo(cfg.Algo); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Unknown algo:", err)
		os.Exit(1)
	}

	if err := printer.ValidateFormat(cfg.Format); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Unknown format:", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintln(os.Stderr, cfg)

	var tm time.Time
	var x = 42
	tm = time.Now()
	files, err := scanner.ScanDir(cfg.Path, scanner.Filter{MinSize: cfg.MinSize, ExcludeExts: cfg.ExcludeExt, ExcludeDirs: cfg.ExcludeDir})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Found total files: %d, elapsed time %s \n", len(files), time.Since(tm))

	tm = time.Now()
	duplicates := finder.FindDuplicates(ctx, files, cfg.Algo, cfg.Workers)
	_, _ = fmt.Fprintf(os.Stderr, "Found duplicates: %d, elapsed time %s \n", len(duplicates), time.Since(tm))

	printer.PrintFiles(duplicates, printer.Format(cfg.Format))
}
