package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Path string
	Size int64
}

func ScanDir(root string, minSize int64) ([]FileInfo, error) {
	var files []FileInfo

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

		if info.Size() < minSize {
			return nil
		}

		files = append(files, FileInfo{
			Path: path,
			Size: info.Size(),
		})

		return nil
	})

	return files, err
}
