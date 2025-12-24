package model

type FileInfo struct {
	Path string
	Size int64
}

type Result struct {
	Hash string
	File FileInfo
	Err  error
}
