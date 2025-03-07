package main

// import serial
import (
	"mputils/pyboard"
	"time"
)

func main() {

	board := pyboard.NewPyboard("COM6")

	board.EnterRawREPL()

	time.Sleep(1 * time.Second)

	board.ExitRawREPL()

}
