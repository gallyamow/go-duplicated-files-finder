package finder_test

import (
	"context"
	"os"
	"testing"

	"github.com/gallyamow/go-duplicated-files-finder/internal/finder"
	"github.com/gallyamow/go-duplicated-files-finder/internal/model"
)

func createTempFile(t *testing.T, dir, content string) model.FileInfo {
	t.Helper()

	f, err := os.CreateTemp(dir, "file_*")
	if err != nil {
		t.Fatalf("failed to create tmp file: %v", err)
	}

	_, err = f.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("failed to close tmp file: %v", err)
	}

	info, err := os.Stat(f.Name())
	if err != nil {
		t.Fatalf("failed to get FileInfo: %v", err)
	}

	return model.FileInfo{
		Path: f.Name(),
		Size: info.Size(),
	}
}

func TestFindDuplicates(t *testing.T) {
	dir := t.TempDir()

	// создаем файлы
	file1 := createTempFile(t, dir, "hello world")
	file2 := createTempFile(t, dir, "hello world")
	file3 := createTempFile(t, dir, "different content")
	file4 := createTempFile(t, dir, "hello world")

	files := []model.FileInfo{file1, file2, file3, file4}

	ctx := context.Background()
	duplicates := finder.FindDuplicates(ctx, files, "md5", 2)

	want := map[string]bool{
		file1.Path: true,
		file2.Path: true,
		file4.Path: true,
	}

	if len(duplicates) != len(want) {
		t.Fatalf("want %d, got %d", len(want), len(duplicates))
	}

	for _, f := range duplicates {
		if !want[f.Path] {
			t.Errorf("got unexpected file: %s", f.Path)
		}
	}
}

func TestFindDuplicates_NoDuplicates(t *testing.T) {
	dir := t.TempDir()

	file1 := createTempFile(t, dir, "a")
	file2 := createTempFile(t, dir, "b")
	file3 := createTempFile(t, dir, "c")

	files := []model.FileInfo{file1, file2, file3}

	ctx := context.Background()
	duplicates := finder.FindDuplicates(ctx, files, "sha256", 2)

	if len(duplicates) != 0 {
		t.Errorf("want 0, got %d", len(duplicates))
	}
}
