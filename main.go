package main

import (
	"mputils/pyboard"
)

func main() {

	board := pyboard.NewPyboard("COM6")

	test, err := board.Exec("print(\"Hello world\")")
	// print the first char as a byte
	println(test, err)

	//println(board.FS.ReadFile("word_clock.py"))

}
