package pyboard

import (
	"bytes"
	"strings"
	"time"

	"go.bug.st/serial"
)

// Pyboard is a struct that represents a pyboard.
type Pyboard struct {
	Port      string
	Serial    serial.Port
	inRawMode bool
	FS        *PyFileSystem
}

// NewPyboard creates a new Pyboard.
func NewPyboard(port string) *Pyboard {
	board := &Pyboard{
		Port:      port,
		inRawMode: false,
		FS:        NewPyFileSystem(),
	}
	board.FS.pyboard = board // set the pyboard for the filesystem

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

// ReadUntil reads from the pyboard's serial connection until a specified prompt is found or timeout occurs.
// It takes a prompt string and optional arguments for max length and timeout duration.
//
// Parameters:
//   - prompt: The string to search for in the incoming data. If empty, reads until buffer is empty
//   - args[0] (optional): Maximum length of returned string. If -1 or not provided, no length limit
//   - args[1] (optional): Timeout in seconds. Default is 10 seconds
//
// Returns:
//   - string: The data read from the serial connection
//   - bool: True if prompt was found or max length reached, false if timeout occurred
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

		//remaining := time.Duration(timeout)*time.Second - time.Since(start)
		p.Serial.SetReadTimeout(25 * time.Millisecond)

		tmpBuffer := make([]byte, 256)
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
			if n > 0 {
				for i := 0; i < n; i++ {
					buffer.WriteByte(tmpBuffer[i])
				}
			}
			if n == 0 && buffer.Len() > 0 {
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

	// send ctrl-b incase we are already in raw repl
	p.Serial.Write([]byte{0x02})

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

	//println("Entered raw repl")
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

	//println("Exited raw repl")
	p.inRawMode = false
	return true

}

func (p *Pyboard) Exec(code string) (string, bool) {

	p.EnterRawREPL()

	p.Serial.Write([]byte(code))

	// write ctrl-D to end the raw repl
	p.Serial.Write([]byte{0x04})

	// read the output
	if _, succ := p.ReadUntil("OK", 2); !succ {
		println("Error executing the code")
	}

	out, _ := p.ReadUntil("", -1, 3, 1)

	endRemove := 3
	if out[0] == 0x04 {
		out = out[1:]
		endRemove = 2
	}

	p.ExitRawREPL()

	var outValue string
	if len(out) > endRemove {
		outValue = strings.TrimSpace(out[:len(out)-endRemove])
	}

	//println("\"" + outValue + "\"")

	return outValue, endRemove != 3
}
