package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

func test1() {
	test2()
}

func test2() {
	test3()
}

func test3() {
	fmt.Printf("%s", debug.Stack())
	debug.PrintStack()
}

func main() {
	for {
		test1()
		time.Sleep(1 * time.Second)
	}

}
