package main

import (
	"hangman-web/src"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go src.StartServer(&wg)
	wg.Wait()
}
