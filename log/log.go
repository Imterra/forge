package log

import (
	"fmt"
	"os"
)

type Exit struct{ Code int }

func HandleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok {
			os.Exit(exit.Code)
		}
		panic(e)
	}
}

func Error(msg string) {
	fmt.Printf("\033[1;31m[ERROR]\033[0m %s\n", msg)
	panic(Exit{1})
}

func Warn(msg string) {
	fmt.Printf("\033[1;33m[WARNING]\033[0m %s\n", msg)
}

func Succ(msg string) {
	fmt.Printf("\033[1;32m[DONE]\033[0m %s\n", msg)
}
