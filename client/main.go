package main

import (
	"../log"
	"./actions"
	"./target"
	"./util"
	"./worker"
	"flag"
	"fmt"
	"os"
	"strings"
)

const ROOT_DEFAULT = "~/.forge"

type workers []*worker.Worker

func (w *workers) String() string {
	ws := make([]string, len(*w))
	for i := range *w {
		ws[i] = (*w)[i].Addr
	}
	return strings.Join(ws, ",")
}

func (w *workers) Set(value string) error {
	for _, w_host := range strings.Split(value, ",") {
		worker, err := worker.GetWorker(w_host)
		if err != nil {
			return err
		}
		*w = append(*w, worker)
	}
	return nil
}

func main() {

	root_flag := flag.String("root", "",
		"Specify root directory for Forge packages.")

	var workers_flag workers
	// TODO: Add locally-running worker automatically.
	flag.Var(
		&workers_flag, "worker",
		"comma-separated list of worker addresses (host:port)")

	flag.Parse()

	var forge_root *string
	forge_root = new(string)
	*forge_root = ROOT_DEFAULT
	root_env := os.Getenv("FORGE_ROOT")
	if root_env != "" {
		forge_root = &root_env
	}
	if *root_flag != "" {
		forge_root = root_flag
	}

	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "\n\nNo target specified.\n\n")
		fmt.Fprintf(os.Stderr, "usage: %s target...\n\n", os.Args[0])
		os.Exit(1)
	}
	targets := flag.Args()
	wd, _ := os.Getwd()

	notifier := make(chan *actions.File, len(targets))
	var conf util.Config

	for i := range targets {
		target_name := targets[i]
		requested_target := target.MakeTarget(target_name, *forge_root, wd)
		requested_file := requested_target.GetOutputFile()
		go actions.MakeFile(requested_file, &conf, notifier)
	}

	for _ = range targets {
		<-notifier
	}
	// TODO: Write metadata for all files.

	log.Succ("Everything done")
}
