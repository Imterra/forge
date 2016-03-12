package main

import (
	"./target"
	"./util"
	//	"flag"
	//	"os"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	hello_c := target.LibCTarget{
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
}
