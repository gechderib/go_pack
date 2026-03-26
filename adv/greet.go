package main

import (
	"fmt"
	"sync"
)

func sayHello(name string, wg *sync.WaitGroup) string {
	defer wg.Done()
	fmt.Println("Saying hello to", name)
	return "Hello, " + name + "!"
}
