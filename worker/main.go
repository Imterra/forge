package main

import (
	"./executor"
	"fmt"
	"os"
)

func main() {
	t1 := executor.CompileCTask{Filename: "auvm.c"}
	t2 := executor.CompileCTask{Filename: "auvmlib.c"}
	t3 := executor.CompileCTask{Filename: "init.c"}
	t4 := executor.CompileCTask{Filename: "ins.c"}
	t5 := executor.CompileCTask{Filename: "intable.c"}
	t6 := executor.CompileCTask{Filename: "object.c"}
	t7 := executor.CompileCTask{Filename: "parse.c"}
	t8 := executor.CompileCTask{Filename: "stack.c"}
	t9 := executor.CompileCTask{Filename: "util.c"}
	t10 := executor.CompileCTask{Filename: "lib/io.c"}

	to_archive := []string{"auvmlib.o", "init.o", "ins.o", "intable.o", "object.o", "parse.o", "stack.o", "util.o", "lib/io.o"}
	tar := executor.ArTask{Name: "libauvm", Inputs: to_archive}
	to_link := []string{"auvm.o", "libauvm.a"}
	tlink := executor.LinkTask{Name: "auvm", Inputs: to_link}
	q := executor.QuitTask{}

	tq := make(chan executor.Task, 100)
	rq := make(chan executor.Response, 100)

	//config := executor.Config{Pkgroot: "/home/pepol/src/auvm", Outdir: "/tmp"}
	config := executor.Config{Pkgroot: "/home/pepol/src/auvm", Outdir: "/tmp/auvm"}

	go executor.Executor(0, tq, rq, &config)
	go executor.Executor(1, tq, rq, &config)
	go executor.Executor(2, tq, rq, &config)

	tq <- &t1
	tq <- &t2
	tq <- &t3
	tq <- &t4
	tq <- &t5
	tq <- &t6
	tq <- &t7
	tq <- &t8
	tq <- &t9
	tq <- &t10

	for i := 0; i < 10; i++ {
		r := <-rq
		if r.Code == executor.ERROR {
			fmt.Printf("\033[1;31m[ERROR] when generating %v, following happened:\033[0m %v\n", r.Outfile, r.Output)
			tq <- &q
			<-rq
			tq <- &q
			<-rq
			tq <- &q
			<-rq
			fmt.Println("Aborting...")
			os.Exit(1)
		} else if r.Code == executor.DONE {
			fmt.Printf("\033[0;32m[DONE] %v\033[0m\n", r.Outfile)
		}
	}

	tq <- &tar
	r := <-rq
	if r.Code == executor.ERROR {
		fmt.Printf("\033[1;31m[ERROR] when generating %v, following happened:\033[0m %v\n", r.Outfile, r.Output)
		tq <- &q
		<-rq
		tq <- &q
		<-rq
		tq <- &q
		<-rq
		fmt.Println("Aborting...")
		os.Exit(1)
	} else if r.Code == executor.DONE {
		fmt.Printf("\033[0;32m[DONE] %v\033[0m\n", r.Outfile)
	}

	tq <- &tlink

	tq <- &q
	tq <- &q
	tq <- &q

	done := make([]bool, 3)
	done[0] = false
	done[1] = false
	done[2] = false

	for {
		r := <-rq
		if r.Code == executor.QUIT {
			done[r.Id] = true
			tester := true
			for i := range done {
				tester = tester && done[i]
			}
			if tester {
				break
			}
		} else if r.Code == executor.ERROR {
			fmt.Printf("\033[1;31m[ERROR]when generating %v, following happened:\033[0m %v\n", r.Outfile, r.Output)
		} else if r.Code == executor.DONE {
			fmt.Printf("\033[0;32m[DONE] %v\033[0m\n", r.Outfile)
		}
	}

	fmt.Println("We're done here!")
}
