package scanner

import (
	"fmt"
	"github.com/gallyamow/go-duplicated-files-finder/internal/model"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

type Filter struct {
	MinSize     int64
	ExcludedExt []string
	IncludedExt []string
}

// ScanDir scans the directory and returns a list of files.
// IO-bound task, no reason to use goroutines.
func ScanDir(root string, filter Filter) ([]model.FileInfo, error) {
	var files []model.FileInfo

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		// symlinks?

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if isExcluded(info, filter) {
			return nil
		}

		files = append(files, model.FileInfo{
			Path: path,
			Size: info.Size(),
		})

		return nil
	})

	return files, err
}

func isExcluded(info fs.FileInfo, filter Filter) bool {
	// size
	if filter.MinSize != 0 {
		if info.Size() < filter.MinSize {
			return true
		}
	}

	// ext
	if filter.ExcludedExt != nil || filter.IncludedExt != nil {
		ext := filepath.Ext(info.Name())

		if filter.ExcludedExt != nil {
			if slices.Contains(filter.ExcludedExt, ext) {
				return true
			}
		}

		if filter.IncludedExt != nil {
			if !slices.Contains(filter.IncludedExt, ext) {
				return true
			}
		}
	}

	return false
}
