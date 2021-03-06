package util

import (
	"../worker"
	"fmt"
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

var Exiter func() = func() {}

func WriteMetadata(rootdir string) error {
	const script = "cd %s && find -type f -a \\! -path '*/.git/*' -a \\! -path '*/.metadata/*' -exec sh -c \"dirname {} | xargs -I[] mkdir -p .metadata/[] && sha512sum {} | cut -d ' ' -f 1 > .metadata/{}\" \\;"

	return exec.Command("sh", "-c", fmt.Sprintf(script, rootdir)).Start()
}
