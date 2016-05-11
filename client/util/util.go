package util

import (
	"../worker"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Rootdir string
	Workers []*worker.Worker
}

func NormalizePath(path string) (string, string) {
	clean_path := filepath.Clean(path)
	full_path, _ := filepath.Abs(clean_path)
	return filepath.Base(full_path), full_path
}

func GetFullPath(abs_filename, rootdir string) string {
	return filepath.Join(rootdir, strings.TrimPrefix(abs_filename, "//"))
}

func CleanupChild(cmd *exec.Cmd) {
	cmd.Process.Signal(os.Interrupt)
	cmd.Wait()
}
