package main

import (
	"./objectstorage"
	"./target"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	//	"path/filepath"
	//	"./objectstorage"
)

const ROOT_DEFAULT = "~/.forge"

func main() {

	root_flag := flag.String("root", "",
		"Specify root directory for Forge packages.")

	//simulate_flag := flag.Bool("simulate", false,
	//	"Display all actions instead of performing them.")

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
	s := objectstorage.FileStorage{Root: "/home/pepol/src/forge_binfiles"}

	for i := 0; i < len(targets); i++ {
		target_name := targets[i]
		requested_target := target.MakeTarget(target_name, *forge_root, wd)
		// target.ParseFile(build_file, target_name)
		fmt.Printf("%v\n", requested_target)
		spew.Dump(requested_target.GetAction(&s))
		fmt.Printf("\n\nAction necessary: %v\nChecksum: %v\n\n", requested_target.GetAction(&s).IsRequired(), requested_target.GetAction(&s).GetOutFilePath())
	}

}

/*	hello_c := target.LibCTarget{
		Name:         "libhello",
		Sources:      []string{"/root/hello.c"},
		Resources:    []string{},
		Dependencies: []target.Target{}}
	hello_out := target.AppCTarget{
		Name:         "hello",
		Sources:      []string{},
		Resources:    []string{},
		Dependencies: []target.Target{&hello_c}}

	a := hello_out.GetAction()

	spew.Dump(a)

	util.PrintAllActions(a)
*/
