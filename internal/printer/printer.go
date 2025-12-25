package printer

import (
	"fmt"
	"github.com/gallyamow/go-duplicated-files-finder/internal/model"
	"os"
	"slices"
)

type Format string

const (
	FormatPlain Format = "plain"
	FormatPaths Format = "paths"
)

func PrintFiles(files []model.FileInfo, format Format) {
	switch format {
	case FormatPlain:
		printFilesPlain(files)
	case FormatPaths:
		printFilesPaths(files)
	}
}

func printFilesPlain(files []model.FileInfo) {
	for _, f := range files {
		if f.Err != nil {
			continue // ошибки не печатаем в stdout
		}
		_, _ = fmt.Fprintf(os.Stdout, "%s\t%d\t%s\n", f.Hash, f.Size, f.Path)
	}
}

func printFilesPaths(files []model.FileInfo) {
	for _, f := range files {
		if f.Err != nil {
			continue // ошибки не печатаем в stdout
		}
		_, _ = fmt.Fprintln(os.Stdout, f.Path)
	}
}

func ValidateFormat(format string) error {
	if slices.Contains([]Format{FormatPlain, FormatPaths}, Format(format)) {
		return nil
	}
	return fmt.Errorf("unknown format: %s", format)
}
