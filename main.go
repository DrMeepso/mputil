package main

import (
	"mputil/pyboard"
	"mputil/tools"
	"os"
	"strings"

	"go.bug.st/serial"
)

func main() {

	// get cli arguments
	args := os.Args[1:]
	if len(args) == 0 {
		Usage()
		return
	}

	// parse arguments, look for -args then command
	var device string
	var command string

	for i := 0; i < len(args); i++ {
		if args[i] == "-d" || args[i] == "--device" {
			device = args[i+1]
			i++
		} else {
			command = args[i]
			break
		}
	}

	var selectedBoard *pyboard.Pyboard
	if device != "" {
		selectedBoard = pyboard.NewPyboard(device)
		defer selectedBoard.Close()
	}

	switch command {

	case "list":
		ListDevices()

	case "exec":
		tools.Tool_Exec(args, selectedBoard)
		return

	case "repl":
		tools.Tool_Repl(args, selectedBoard)
		return

	default:
		println("Unknown command")
		Usage()

	}

}

func Usage() {

	split := strings.Split(os.Args[0], "\\")
	exeName := split[len(split)-1]

	println("Usage:")
	println("  " + exeName + " [options] <command>")
	println("")
	println("Options:")
	println("  -d, --device <comport>  Specify the device comport")
	println("")
	println("Tools:")
	println("  list                    List available comports")
	println("  dump <local folder>     Dump the pyboard filesystem to the local folder")
	println("  sync <local folder>     Sync the local folder to the pyboard")
	println("  exec                    Execute python code on the pyboard")
	println("  repl                    Start a python repl on the pyboard")

}

func ListDevices() {

	// list available ports
	println("Available comports:")
	allPosts, _ := serial.GetPortsList()
	for _, port := range allPosts {
		println("  " + port)
	}

}
