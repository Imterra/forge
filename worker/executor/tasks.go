package executor

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type CompileCTask struct {
	Filename string
}

func (t *CompileCTask) GetType() string {
	return "compile_c"
}

func (t *CompileCTask) GetOutfile() string {
	if !strings.HasSuffix(t.Filename, ".c") {
		return t.Filename + ".o"
	}
	return strings.TrimSuffix(t.Filename, ".c") + ".o"
}

func (t *CompileCTask) GetInputs() []string {
	ins := make([]string, 1)
	ins[0] = t.Filename
	return ins
}

func (t *CompileCTask) Execute(config *Config) (string, error) {
	outfile_path := filepath.Join(config.Outdir, t.GetOutfile())
	infile_path := filepath.Join(config.Pkgroot, t.Filename)

	output, err := exec.Command("gcc", "-std=c99", "-c", "-o", outfile_path, infile_path).CombinedOutput()

	return string(output), err
}

type ArTask struct {
	Name   string
	Inputs []string
}

func (t *ArTask) GetType() string {
	return "ar"
}

func (t *ArTask) GetOutfile() string {
	return t.Name + ".a"
}

func (t *ArTask) GetInputs() []string {
	return t.Inputs
}

func (t *ArTask) Execute(config *Config) (string, error) {
	outfile_path := filepath.Join(config.Outdir, t.GetOutfile())

	args := make([]string, len(t.Inputs)+2)
	for i := range t.Inputs {
		args[i+2] = filepath.Join(config.Outdir, t.Inputs[i])
	}

	args[0] = "rcs"
	args[1] = outfile_path

	output, err := exec.Command("ar", args...).CombinedOutput()

	return string(output), err
}

type LinkTask struct {
	Name   string
	Inputs []string
}

func (t *LinkTask) GetType() string {
	return "link"
}

func (t *LinkTask) GetOutfile() string {
	return t.Name
}

func (t *LinkTask) GetInputs() []string {
	return t.Inputs
}

func (t *LinkTask) Execute(config *Config) (string, error) {
	outfile_path := filepath.Join(config.Outdir, t.GetOutfile())

	args := make([]string, len(t.Inputs)+2)
	for i := range t.Inputs {
		args[i+2] = filepath.Join(config.Outdir, t.Inputs[i])
	}

	args[0] = "-o"
	args[1] = outfile_path

	output, err := exec.Command("gcc", args...).CombinedOutput()

	return string(output), err
}
