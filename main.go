package main

import (
	"mputils/pyboard"
)

func main() {

	board := pyboard.NewPyboard("COM6")

	println(board.FS.ReadFile("word_clock.py"))

}
