package config

import (
	"errors"
	"os"
	"reflect"
	"runtime"
	"testing"
)

func TestParseFlags(t *testing.T) {
	t.Run("no path", func(t *testing.T) {
		resetFlags()

		os.Args = []string{"cmd"}

		_, err := ParseFlags()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrorPathRequired) {
			t.Fatalf("got %q, want %q", err, ErrorPathRequired)
		}
	})

	t.Run("use defaults", func(t *testing.T) {
		resetFlags()

		os.Args = []string{
			"cmd",
			"/test",
		}

		cfg, err := ParseFlags()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ExcludeExt != nil {
			t.Fatalf("got %v, want %v", cfg.ExcludeExt, nil)
		}

		if cfg.ExcludeDir != nil {
			t.Fatalf("got %v, want %v", cfg.ExcludeDir, nil)
		}

		if cfg.MinSize != 0 {
			t.Fatalf("got %v, want %q", cfg.MinSize, 0)
		}

		if cfg.Format != "plain" {
			t.Fatalf("got %v, want %v", cfg.Format, "plain")
		}

		if cfg.Algo != "md5" {
			t.Fatalf("got %v, want %v", cfg.Algo, "md5")
		}

		if cfg.Workers != runtime.GOMAXPROCS(0)*2 {
			t.Fatalf("got %q, want %q", cfg.Workers, runtime.GOMAXPROCS(0)*2)
		}
	})

	t.Run("correct parsing short form", func(t *testing.T) {
		resetFlags()

		os.Args = []string{
			"cmd",
			"-min-size=3Mb",
			"-exclude-ext=.jpg,.png",
			"-exclude-dir=logs,cache",
			"-algo=sha256",
			"-workers=2",
			"-format=plain",
			"/test",
		}

		got, err := ParseFlags()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := &Config{
			Path:       "/test",
			MinSize:    3 * 1024 * 1024,
			ExcludeExt: []string{".jpg", ".png"},
			ExcludeDir: []string{"logs", "cache"},
			Algo:       "sha256",
			Workers:    2,
			Format:     "plain",
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("correct parsing long form", func(t *testing.T) {
		resetFlags()

		os.Args = []string{
			"cmd",
			"--min-size=3Mb",
			"--exclude-ext=.jpg,.png",
			"--exclude-dir=logs,cache",
			"--algo=sha256",
			"--workers=2",
			"--format=plain",
			"/test",
		}

		got, err := ParseFlags()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := &Config{
			Path:       "/test",
			MinSize:    3 * 1024 * 1024,
			ExcludeExt: []string{".jpg", ".png"},
			ExcludeDir: []string{"logs", "cache"},
			Algo:       "sha256",
			Workers:    2,
			Format:     "plain",
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})
}

func TestParseStrArray(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty", "", nil},
		{"single", "jpg", []string{"jpg"}},
		{"multiple", "jpg,png,gif", []string{"jpg", "png", "gif"}},
		{"with spaces", " jpg, png , gif ", []string{"jpg", "png", "gif"}},
		{"with empty values", "jpg,,png,", []string{"jpg", "png"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseStrArray(tt.input)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"bytes", "10B", 10, false},
		{"kilobytes", "1KB", 1024, false},
		{"megabytes", "2MB", 2 * 1024 * 1024, false},
		{"gigabytes", "1GB", 1024 * 1024 * 1024, false},
		{"float value", "1.5KB", 1536, false},
		{"invalid number", "abcKB", 0, true},
		{"unknown unit", "10TB", 0, true},
		{"empty", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSize(tt.input)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
