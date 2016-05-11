package tasks

import (
	"../../proto"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Root string
}

type Task struct {
	Semaphore chan int
	Config    *Config
}

func (t *Task) CompileC(args *proto.Args, resp *proto.Response) error {
	t.Semaphore <- 1
	defer func() { <-t.Semaphore }()

	var outfilename string
	if len(args.Inputs) < 1 {
		return errors.New("compile requires at least one file to be compiled")
	}

	if !strings.HasSuffix(args.Inputs[0].Filename, ".c") {
		outfilename = args.Inputs[0].Filename + ".o"
	} else {
		outfilename = strings.TrimSuffix(args.Inputs[0].Filename, ".c") + ".o"
	}

	outfile_path := filepath.Join(t.Config.Root, outfilename)
	infile_path := filepath.Join(t.Config.Root, args.Inputs[0].Filename)

	output, err := exec.Command("gcc", "-std=c99", "-c", "-o",
		outfile_path, infile_path).CombinedOutput()

	if err != nil {
		return err
	}

	err = prepareResponse(outfile_path, output, args.SendContent, resp)
	return err
}

func (t *Task) ArLink(args *proto.Args, resp *proto.Response) error {
	t.Semaphore <- 1
	defer func() { <-t.Semaphore }()

	var outfilename string
	if len(args.Inputs) < 1 {
		return errors.New("library linking requires at least one file")
	}

	outfilename = args.Name + ".a"
	outdir := t.Config.Root

	outfile_path := filepath.Join(outdir, outfilename)

	ar_args := processInputs(args.Inputs, outdir, []string{"rcs", outfile_path})

	output, err := exec.Command("ar", ar_args...).CombinedOutput()
	if err != nil {
		return err
	}

	err = prepareResponse(outfile_path, output, args.SendContent, resp)
	return err
}

func (t *Task) LdLink(args *proto.Args, resp *proto.Response) error {
	t.Semaphore <- 1
	defer func() { <-t.Semaphore }()

	var outfilename string
	if len(args.Inputs) < 1 {
		return errors.New("application linking requires at least one file")
	}

	outfilename = args.Name
	outdir := t.Config.Root

	outfile_path := filepath.Join(outdir, outfilename)

	ld_args := processInputs(args.Inputs, outdir, []string{"-o", outfile_path})

	output, err := exec.Command("gcc", ld_args...).CombinedOutput()
	if err != nil {
		return err
	}

	err = prepareResponse(outfile_path, output, args.SendContent, resp)
	return err
}
