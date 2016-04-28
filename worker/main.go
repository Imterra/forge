package main

import (
	"./executor"
	"fmt"
	"runtime"
)

func main() {
	a := make([]string, 0)
	s := make([]string, 1)
	s[0] = "foo.c"
	t := executor.Task{
		Type:        "compile_c",
		OutFilename: "foo.o",
		SrcArgs:     s,
		ResArgs:     a,
		DepArgs:     a}
	q := executor.Task{
		Type:        "stop",
		OutFilename: "",
		SrcArgs:     nil,
		ResArgs:     nil,
		DepArgs:     nil}

	tq := make(chan executor.Task, runtime.NumCPU())
	rq := make(chan executor.Response, runtime.NumCPU())

	go executor.Executor(1, tq, rq)
	go executor.Executor(2, tq, rq)

	fmt.Println("Sending T")
	tq <- t
	tq <- t
	tq <- t
	tq <- t
	tq <- q

	fmt.Println("Waiting for T to be done")

	for {
		r := <-rq
		fmt.Printf("GOT: %v\n", r)
		if r.Code == executor.QUIT {
			break
		}
	}

	fmt.Println("We're done here!")
}
