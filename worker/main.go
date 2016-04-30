package main

import (
	"./tasks"
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

	port_flag := flag.Int("port", 1103, "Specify port number to listen on.")
	root_flag := flag.String("root", "",
		"Specify root working directory for Forge.")
	flag.Parse()

	var port int
	port = PORT_DEFAULT
	if *port_flag > 0 && *port_flag < 65535 {
		port = *port_flag
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

	sem := make(chan int, runtime.NumCPU())
	config := tasks.Config{Root: *forge_root}

	task := tasks.Task{Semaphore: sem, Config: &config}
	rpc.Register(&task)

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
