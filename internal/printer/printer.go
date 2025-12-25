package printer

import (
	"fmt"
	"os"
	"slices"
	"sort"

	"github.com/gallyamow/go-duplicated-files-finder/internal/model"
)

type Format string

const (
	FormatPlain Format = "plain"
	FormatPaths Format = "paths"
)

func ValidateFormat(format string) error {
	if slices.Contains([]Format{FormatPlain, FormatPaths}, Format(format)) {
		return nil
	}
	return fmt.Errorf("unknown format: %s", format)
}

func PrintFiles(files []model.FileInfo, format Format) {
	switch format {
	case FormatPlain:
		printFilesPlain(files)
	case FormatPaths:
		printFilesPaths(files)
	}
}

func printFilesPlain(files []model.FileInfo) {
	sortedFiles := sorted(files)

	for _, f := range sortedFiles {
		if f.Err != nil {
			continue // ошибки не печатаем в stdout
		}
		_, _ = fmt.Fprintf(os.Stdout, "%s\t%s\t%s\n", f.Hash, formatFileSize(f.Size), f.Path)
	}
}

func printFilesPaths(files []model.FileInfo) {
	sortedFiles := sorted(files)

	for _, f := range sortedFiles {
		if f.Err != nil {
			continue // ошибки не печатаем в stdout
		}
		_, _ = fmt.Fprintln(os.Stdout, f.Path)
	}
}

func sorted(files []model.FileInfo) []model.FileInfo {
	res := append([]model.FileInfo{}, files...)

	sort.Slice(res, func(i, j int) bool {
		return res[i].Size > res[j].Size
	})

	return res
}

func formatFileSize(bytes int64) string {
	const unit = 1000
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB",
		float64(bytes)/float64(div),
		"KMGTPE"[exp],
	)
}
