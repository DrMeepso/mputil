package pyboard

import (
	"bytes"
	"time"

	"go.bug.st/serial"
)

// Pyboard is a struct that represents a pyboard.
type Pyboard struct {
	Port      string
	Serial    serial.Port
	inRawMode bool
}

// NewPyboard creates a new Pyboard.
func NewPyboard(port string) *Pyboard {
	board := &Pyboard{
		Port:      port,
		inRawMode: false,
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

func (p *Pyboard) Close() {
	p.Serial.Close()
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

	timeout := 10 // default timeout is 10 seconds
	// the 2nd argument is the timeout
	if len(args) > 1 {
		timeout = args[1]
	}

	defer func() {
		p.Serial.SetReadTimeout(serial.NoTimeout) // reset the timeout
	}()

	start := time.Now()
	for {

		if time.Since(start).Seconds() > float64(timeout) {
			return buffer.String(), false
		}

		tmpBuffer := make([]byte, 1)
		remaining := time.Duration(timeout)*time.Second - time.Since(start)
		p.Serial.SetReadTimeout(remaining)

		n, err := p.Serial.Read(tmpBuffer)
		if err != nil {
			panic(err)
		}

		// loop through the tmpBuffer and append to the buffer, until the prompt is found
		if prompt != "" {
			for i := 0; i < n; i++ {
				buffer.WriteByte(tmpBuffer[i])
				if bytes.Contains(buffer.Bytes(), []byte(prompt)) {
					if maxLength == -1 {
						return buffer.String(), true
					} else if buffer.Len() >= maxLength {
						start := buffer.Len() - maxLength
						return string(buffer.Bytes()[start:]), true
					}
				}
			}
		} else {
			buffer.Write(tmpBuffer)
			// if n == 0, it means that the buffer is empty
			if n == 0 {
				if maxLength == -1 {
					return buffer.String(), true
				} else if buffer.Len() >= maxLength {
					start := buffer.Len() - maxLength
					return string(buffer.Bytes()[start:]), true
				}
			}
		}
	}
}

func (p *Pyboard) EnterRawREPL() bool {

	if p.inRawMode {
		println("Already in raw repl")
		return true
	}

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
	p.inRawMode = true
	return true

}

func (p *Pyboard) ExitRawREPL() bool {

	if !p.inRawMode {
		println("Not in raw repl")
		return true
	}

	// send ctrl-b to exit raw repl
	p.Serial.Write([]byte{0x02})

	if _, found := p.ReadUntil(">>>", 3, 3); !found {
		println("Unable to exit raw repl")
		return false
	}

	println("Exited raw repl")
	return true

}

func (p *Pyboard) Exec(code string) string {

	p.EnterRawREPL()

	p.Serial.Write([]byte(code))

	// write ctrl-D to end the raw repl
	p.Serial.Write([]byte{0x04})

	// read the output
	if _, succ := p.ReadUntil("OK", 2); !succ {
		println("Error executing the code")
	}

	out, _ := p.ReadUntil("", -1, 3)

	p.ExitRawREPL()

	return out[:len(out)-4]
}
