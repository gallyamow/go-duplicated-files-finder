package config

import (
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type Config struct {
	Path       string
	Delete     bool
	MinSize    int64
	ExcludeExt []string
	ExcludeDir []string
	Algo       string
	Workers    int
	Format     string
}

func (c *Config) String() string {
	return fmt.Sprintf("Config { Path: %s, Delete: %t, MinSize: %d, excludedExts: %v, excludedDirs: %v, Algo: %s, Workers: %d, Format: %s }",
		c.Path, c.Delete, c.MinSize, c.ExcludeExt, c.ExcludeDir, c.Algo, c.Workers, c.Format)
}

func ParseFlags() (*Config, error) {
	deleteFlag := flag.Bool("delete", false, "delete duplicate files")
	minSizeStr := flag.String("min-size", "1B", "minimum file size (e.g. 10MB, 500KB)")
	excludeExtStr := flag.String("exclude-ext", "", "comma-separated extensions")
	excludeDirStr := flag.String("exclude-dir", "", "comma-separated directory names")
	algo := flag.String("algo", "md5", "hash algorithm: md5, sha1, sha256")
	workers := flag.Int("workers", runtime.GOMAXPROCS(0)*2, "number of concurrent workers")
	format := flag.String("format", "plain", "output format: plain, paths")

	flag.Parse()

	if flag.NArg() < 1 {
		return nil, fmt.Errorf("path is required")
	}

	path := flag.Arg(0)

	minSize, err := parseSize(*minSizeStr)
	if err != nil {
		return nil, err
	}

	excludeExts := parseStrArray(strings.ToLower(*excludeExtStr))
	excludeDirs := parseStrArray(strings.ToLower(*excludeDirStr))

	if *workers <= 0 {
		return nil, fmt.Errorf("workers must be > 0")
	}

	if *deleteFlag && *workers > 1 {
		// TODO
		return nil, fmt.Errorf("cannot delete files with multiple workers")
	}

	return &Config{
		Path:       path,
		Delete:     *deleteFlag,
		MinSize:    minSize,
		ExcludeExt: excludeExts,
		ExcludeDir: excludeDirs,
		Algo:       *algo,
		Workers:    *workers,
		Format:     *format,
	}, nil
}

func parseStrArray(s string) []string {
	if s == "" {
		return nil
	}

	var res []string
	for _, v := range strings.Split(s, ",") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		res = append(res, v)
	}
	return res
}

func parseSize(s string) (int64, error) {
	s = strings.ToUpper(strings.TrimSpace(s))

	multipliers := []struct {
		unit string
		mul  int64
	}{
		{"GB", 1024 * 1024 * 1024},
		{"MB", 1024 * 1024},
		{"KB", 1024},
		{"B", 1},
	}

	for _, it := range multipliers {
		if strings.HasSuffix(s, it.unit) {
			value := strings.TrimSuffix(s, it.unit)
			n, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid size: %s", s)
			}
			return int64(n * float64(it.mul)), nil
		}
	}

	return 0, fmt.Errorf("unknown size unit: %s", s)
}
