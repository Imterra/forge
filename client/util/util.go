package util

import (
	"net/rpc"
	"path/filepath"
	"strings"
)

type Config struct {
	Request bool
	Client  *rpc.Client
	Monitor chan *rpc.Call
}

func NormalizePath(path string) (string, string) {
	clean_path := filepath.Clean(path)
	full_path, _ := filepath.Abs(clean_path)
	return filepath.Base(full_path), full_path
}

func GetFullPath(abs_filename, rootdir string) string {
	return filepath.Join(rootdir, strings.TrimPrefix(abs_filename, "//"))
}
