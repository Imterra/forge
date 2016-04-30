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
	if len(args.Inputs) != 1 {
		return errors.New("Compile requires one file to be compiled!")
	}

	if !strings.HasSuffix(args.Inputs[0], ".c") {
		outfilename = args.Inputs[0] + ".o"
	} else {
		outfilename = strings.TrimSuffix(args.Inputs[0], ".c") + ".o"
	}

	outfile_path := filepath.Join(t.Config.Root, "BINFILES", outfilename)
	infile_path := filepath.Join(t.Config.Root, args.Inputs[0])

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
		return errors.New("Library linking requires at least one file!")
	}

	outfilename = args.Name + ".a"
	outdir := filepath.Join(t.Config.Root, "BINFILES")

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
		return errors.New("Application linking required at least one file!")
	}

	outfilename = args.Name
	outdir := filepath.Join(t.Config.Root, "BINFILES")

	outfile_path := filepath.Join(outdir, outfilename)

	ld_args := processInputs(args.Inputs, outdir, []string{"-o", outfile_path})

	output, err := exec.Command("gcc", ld_args...).CombinedOutput()
	if err != nil {
		return err
	}

	err = prepareResponse(outfile_path, output, args.SendContent, resp)
	return err
}
