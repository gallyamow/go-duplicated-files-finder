package scanner

import (
	"fmt"
	"github.com/gallyamow/go-duplicated-files-finder/internal/model"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Filter struct {
	MinSize     int64
	ExcludeExts []string
	ExcludeDirs []string
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
			if isDirExcluded(d, filter) {
				return fs.SkipDir
			}
			return nil
		}

		// symlinks?

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if isFileExcluded(info, filter) {
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

func isFileExcluded(info fs.FileInfo, filter Filter) bool {
	// size
	if filter.MinSize != 0 {
		if info.Size() < filter.MinSize {
			return true
		}
	}

	// ext
	if len(filter.ExcludeExts) > 0 {
		if slices.Contains(filter.ExcludeExts, strings.ToLower(filepath.Ext(info.Name()))) {
			return true
		}
	}

	return false
}

func isDirExcluded(d fs.DirEntry, filter Filter) bool {
	if len(filter.ExcludeDirs) > 0 {
		if slices.Contains(filter.ExcludeDirs, strings.ToLower(d.Name())) {
			return true
		}
	}

	return false
}
