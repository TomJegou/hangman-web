package main

import (
	"hangman-web/src"
	"os"
	"strings"
	"sync"
)

func main() {
	/*
		We want that the user can start the executable file whatever his current working directory is
	*/
	executablepath, _ := os.Executable()    // get the absolute path of the executable
	t := strings.Split(executablepath, "/") //split into a slice all the dir's name
	srcpath := "/"                          // this will be the path where we will redirect the program
	for i := 0; i < len(t)-2; i++ {
		srcpath += t[i] + "/"
	}
	srcpath = srcpath[1:]
	os.Chdir(srcpath)     // change the current working directory
	var wg sync.WaitGroup // creating the waitgroup for the goroutine
	wg.Add(1)
	go src.StartServer(&wg)
	wg.Wait() // wait the goroutine server to finish before the end of the main func
}
