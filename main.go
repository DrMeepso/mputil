package main

import (
	"mputil/pyboard"

	"go.bug.st/serial"
)

func main() {

	// list available ports
	allPosts, _ := serial.GetPortsList()
	for _, port := range allPosts {
		println(port)
	}

	board := pyboard.NewPyboard("COM6")

	// write hello.txt to the pyboard
	board.FS.WriteFile("hello.txt", "Hello, world!")

	println(board.FS.ReadFile("hello.txt"))

	println(board.FS.GetSHA256("hello.txt"))

}
