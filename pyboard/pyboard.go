package pyboard

import (
	"bytes"
	"time"

	"go.bug.st/serial"
)

// Pyboard is a struct that represents a pyboard.
type Pyboard struct {
	Port   string
	Serial serial.Port
}

// NewPyboard creates a new Pyboard.
func NewPyboard(port string) *Pyboard {
	board := &Pyboard{
		Port: port,
	}

	// connect to the pyboard

	// rp2350's connection details
	serial, err := serial.Open(port, &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		panic(err)
	}

	board.Serial = serial

	println("Connected to the pyboard")

	return board
}

// blocks untill it reads the prompt value from the serial port
func (p *Pyboard) ReadUntil(prompt string, args ...int) (string, bool) {
	// read from the pyboard until the prompt is found
	var buffer bytes.Buffer

	// if max is provided, read until max bytes are read
	maxLength := -1
	if len(args) > 0 {
		maxLength = args[0]
	}

	timeout := 10
	// the 2nd argument is the timeout
	if len(args) > 1 {
		timeout = args[1]
	}

	start := time.Now()
	for {

		println("Reading from the pyboard")

		if time.Since(start).Seconds() > float64(timeout) {
			return buffer.String(), false
		}

		tmpBuffer := make([]byte, 1024)
		remaining := time.Duration(timeout)*time.Second - time.Since(start)
		p.Serial.SetReadTimeout(remaining)

		n, err := p.Serial.Read(tmpBuffer)
		if err != nil {
			panic(err)
		}

		// loop through the tmpBuffer and append to the buffer, until the prompt is found
		for i := 0; i < n; i++ {
			buffer.WriteByte(tmpBuffer[i])
			if bytes.Contains(buffer.Bytes(), []byte(prompt)) {
				if maxLength == -1 {
					return buffer.String(), true
				} else if buffer.Len() >= maxLength {
					return string(buffer.Bytes()[:maxLength+1]), true
				}
			}
		}
	}
}

func (p *Pyboard) EnterRawREPL() bool {

	// send ctrl-c to stop the running program
	p.Serial.Write([]byte{0x03})

	// read the output
	if _, found := p.ReadUntil(">>>", 3, 3); !found {
		println("Unable to stop the program")
		return false
	}

	// send ctrl-a to enter raw repl
	p.Serial.Write([]byte{0x01})

	if _, found := p.ReadUntil("CTRL-B to exit", 5, 3); !found {
		println("Unable to enter raw repl")
		return false
	}

	println("Entered raw repl")
	return true

}

func (p *Pyboard) ExitRawREPL() bool {

	// send ctrl-b to exit raw repl
	p.Serial.Write([]byte{0x02})

	if _, found := p.ReadUntil(">>>", 3, 3); !found {
		println("Unable to exit raw repl")
		return false
	}

	println("Exited raw repl")
	return true

}
