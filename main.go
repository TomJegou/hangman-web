package main

import (
	"fmt"
	"os"
	"strings"
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
	fmt.Println(srcpath) // change the current working directory
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go src.StartServer(&wg)
	// wg.Wait()
}
