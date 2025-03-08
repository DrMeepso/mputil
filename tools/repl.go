package tools

import (
	"mputil/pyboard"
	"os"
	"time"
)

func Tool_Repl(args []string, board *pyboard.Pyboard) {

	// send ctrl+c to the board to exit any running code
	board.Serial.Write([]byte{3})

	time.Sleep(10 * time.Millisecond)

	buff := make([]byte, 128)
	board.Serial.Read(buff)

	// create a new go routine to handle the read the serial port
	go func() {
		for {
			buff := make([]byte, 64)
			n, err := board.Serial.Read(buff)
			if err != nil {
				println("Error reading from serial port")
				return
			}
			if n > 0 {
				print(string(buff))
			}
		}
	}()

	println("Entering REPL. Press Ctrl+C to exit.")

	// send ctrl+d to soft reset the board
	board.Serial.Write([]byte{4})

	// read from stdin and write to the serial port
	for {
		buff := make([]byte, 64)
		n, err := os.Stdin.Read(buff)
		if err != nil {
			println("Error reading from stdin")
			return
		}
		if n > 0 {
			_, err := board.Serial.Write(buff[:n])
			if err != nil {
				println("Error writing to serial port")
				return
			}
		}
	}
}
