package hasher

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gallyamow/go-duplicated-files-finder/internal/model"
	"github.com/gallyamow/go-duplicated-files-finder/pkg/workerpool"
	"hash"
	"io"
	"os"
	"slices"
	"sync"
)

func HashFiles(ctx context.Context, files []model.FileInfo, algo string, workers int) []model.Result {
	// hasher is not thread-safe, so each worker must use its own hasher, so we will use sync.Pool in order to
	var hasherPool = sync.Pool{
		New: func() interface{} {
			hasher, err := newHasher(algo)
			if err != nil {
				panic(err)
			}
			return hasher
		},
	}

	handler := func(ctx context.Context, job model.FileInfo) model.Result {
		if ctx.Err() != nil {
			return model.Result{Err: ctx.Err()}
		}

		hasher := hasherPool.Get().(hash.Hash)
		defer func() {
			hasher.Reset()
			hasherPool.Put(hasher)
		}()

		hashStr, err := hashFile(job.Path, algo, hasher)
		if err != nil {
			return model.Result{Err: err}
		}

		return model.Result{
			Hash: hashStr,
			File: job,
		}
	}

	jobCh := make(chan model.FileInfo, len(files))
	go func() {
		defer close(jobCh)

		for _, job := range files {
			select {
			case <-ctx.Done():
				break
			case jobCh <- job:
			}
		}
	}()

	resCh := workerpool.RunWithWorkers[model.FileInfo, model.Result](ctx, jobCh, handler, workers)

	results := make([]model.Result, 0, len(files))
	for res := range resCh {
		results = append(results, res)
	}

	return results
}

func ValidateAlgo(algo string) error {
	if slices.Contains([]string{"md5", "sha1", "sha256"}, algo) {
		return nil
	}
	return fmt.Errorf("unknown algorithm: %s", algo)
}

func hashFile(path string, algo string, hasher hash.Hash) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func newHasher(algo string) (hash.Hash, error) {
	switch algo {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	default:
		return nil, fmt.Errorf("unsupported hash algorithm: %s", algo)
	}
}
