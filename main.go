package main

import (
	"mputils/pyboard"
)

func main() {

	board := pyboard.NewPyboard("COM6")

	println(board.Exec("import os\n\rprint(os.listdir())"))

}
