package main

import (
	"./files"
	"./tasks"
	"./util"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"runtime"
)

const PORT_DEFAULT = 1103
const ROOT_DEFAULT = "~/forge"

func main() {

	port_flag := flag.Int("port", PORT_DEFAULT, "Specify port number to listen on.")
	root_flag := flag.String("root", "",
		"Specify root working directory for Forge.")
	jobs_flag := flag.Int("jobs", runtime.NumCPU(),
		"Specify number of jobs to run simultaneously.")
	flag.Parse()

	var port int
	port = PORT_DEFAULT
	if *port_flag > 0 && *port_flag < 65535 {
		port = *port_flag
	}

	var jobs int
	jobs = runtime.NumCPU()
	if *jobs_flag > 0 {
		jobs = *jobs_flag
	}

	var forge_root *string = new(string)
	*forge_root = ROOT_DEFAULT
	root_env := os.Getenv("FORGE_ROOT")
	if root_env != "" {
		forge_root = &root_env
	}
	if *root_flag != "" {
		forge_root = root_flag
	}

	sem := make(chan int, jobs)
	config := tasks.Config{Root: *forge_root}

	task := tasks.Task{Semaphore: sem, Config: &config}
	file := files.File{Rootdir: *forge_root}
	util := util.Util{Jobs: jobs}
	rpc.Register(&task)
	rpc.Register(&file)
	rpc.Register(&util)

	port_spec := fmt.Sprintf(":%d", port)

	tcp_addr, err := net.ResolveTCPAddr("tcp", port_spec)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcp_addr)
	checkError(err)

	rpc.Accept(listener)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
