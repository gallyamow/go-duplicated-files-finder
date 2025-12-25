package config

import (
	"flag"
	"fmt"
	"github.com/gallyamow/go-duplicated-files-finder/internal/hasher"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Path    string
	Delete  bool
	MinSize int64
	Algo    string
	Workers int
}

func (c *Config) String() string {
	return fmt.Sprintf("Config { Path: %s, Delete: %t, MinSize: %d, Algo: %s, Workers: %d }", c.Path, c.Delete, c.MinSize, c.Algo, c.Workers)
}

var (
	config *Config
	once   sync.Once
)

func ParseFlags() (*Config, error) {
	deleteFlag := flag.Bool("delete", false, "delete duplicate files")
	minSizeStr := flag.String("min-size", "1B", "minimum file size (e.g. 10MB, 500KB)")
	algo := flag.String("algo", "sha256", "hash algorithm: md5, sha1, sha256")
	workers := flag.Int("workers", 4, "number of concurrent workers")

	flag.Parse()

	if flag.NArg() < 1 {
		return nil, fmt.Errorf("path is required")
	}

	path := flag.Arg(0)

	minSize, err := parseSize(*minSizeStr)
	if err != nil {
		return nil, err
	}

	if *workers <= 0 {
		return nil, fmt.Errorf("workers must be > 0")
	}

	if err := hasher.ValidateAlgo(*algo); err != nil {
		return nil, fmt.Errorf("unknown algo %s", *algo)
	}

	if *deleteFlag && *workers > 1 {
		// TODO
		return nil, fmt.Errorf("cannot delete files with multiple workers")
	}

	return &Config{
		Path:    path,
		Delete:  *deleteFlag,
		MinSize: minSize,
		Algo:    *algo,
		Workers: *workers,
	}, nil
}

func parseSize(s string) (int64, error) {
	s = strings.ToUpper(strings.TrimSpace(s))

	multipliers := map[string]int64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
	}

	for unit, mul := range multipliers {
		if strings.HasSuffix(s, unit) {
			value := strings.TrimSuffix(s, unit)
			n, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid size: %s", s)
			}
			return int64(n * float64(mul)), nil
		}
	}

	return 0, fmt.Errorf("unknown size unit: %s", s)
}
