package main

import (
	"./target"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
)

const ROOT_DEFAULT = "~/.forge"

func main() {

	root_flag := flag.String("root", "",
		"Specify root directory for Forge packages.")

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

	for i := 0; i < len(targets); i++ {
		target_name := targets[i]
		requested_target := target.MakeTarget(target_name, *forge_root, wd)
		spew.Dump(requested_target)
		fmt.Println()
		fmt.Println()
		spew.Dump(requested_target.GetOutputFile())
	}
}
