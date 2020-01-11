package main

import (
	"fmt"
	"os"
	"syscall"
)

func Error(e error) {
	fmt.Printf("%s\n", e.Error())
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
}