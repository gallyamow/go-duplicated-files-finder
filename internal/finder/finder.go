package finder

import (
	"context"

	"github.com/gallyamow/go-duplicated-files-finder/internal/hasher"
	"github.com/gallyamow/go-duplicated-files-finder/internal/model"
)

func FindDuplicates(ctx context.Context, files []model.FileInfo, algo string, workers int) []model.FileInfo {
	return findEqualHashed(ctx, findEqualSized(files), algo, workers)
}

func findEqualSized(files []model.FileInfo) []model.FileInfo {
	mp := make(map[int64][]model.FileInfo)
	for _, file := range files {
		mp[file.Size] = append(mp[file.Size], file)
	}

	var res []model.FileInfo
	for _, group := range mp {
		if len(group) == 1 {
			continue
		}
		res = append(res, group...)
	}

	return res
}

func findEqualHashed(ctx context.Context, files []model.FileInfo, algo string, workers int) []model.FileInfo {
	hashedRes := hasher.HashFiles(ctx, files, algo, workers)

	mp := make(map[string][]model.FileInfo)
	for _, res := range hashedRes {
		mp[res.Hash] = append(mp[res.Hash], res)
	}

	var res []model.FileInfo
	for _, group := range mp {
		if len(group) == 1 {
			continue
		}
		res = append(res, group...)
	}

	return res
}
