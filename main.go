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

	board := pyboard.NewPyboard("COM7")

	test, err := board.Exec("print(\"Hello world\")")
	// print the first char as a byte
	println(test, err)

	//println(board.FS.ReadFile("word_clock.py"))

}
