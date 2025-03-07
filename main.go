package main

import (
	"mputils/pyboard"
)

func main() {

	board := pyboard.NewPyboard("COM6")

	println("\"" + board.Exec("print('Hello World')") + "\"")

}
