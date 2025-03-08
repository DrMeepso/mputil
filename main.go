package main

import (
	"bufio"
	"mputil/pyboard"
	"os"
	"strings"

	"go.bug.st/serial"
)

func main() {

	// list available ports
	/*
		allPosts, _ := serial.GetPortsList()
		for _, port := range allPosts {
			println(port)
		}
	*/

	/*
		board := pyboard.NewPyboard("COM6")

		// write hello.txt to the pyboard
		board.FS.WriteFile("hello.txt", "Hello, world!")

		println(board.FS.ReadFile("hello.txt"))

		println(board.FS.GetSHA256("hello.txt"))
	*/

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
		if selectedBoard == nil {
			println("No device selected")
			return
		}
		// prompt the user for python code
		reader := bufio.NewReader(os.Stdin)
		println("Enter the code to eval:")
		print("> ")
		code, _ := reader.ReadString('\n')
		resp, err := selectedBoard.Exec(code)
		if err {
			println("Error executing python code")
			println(resp)
		} else {
			println(resp)
		}
		return
	}

}

func Usage() {
	split := strings.Split(os.Args[0], "\\")
	exeName := split[len(split)-1]

	println("Usage:")
	println("  " + exeName + " [options]")
	println("")
	println("Options:")
	println("  -d, --device <comport>  Specify the device comport")
	println("")
	println("Tools:")
	println("  list                    List available comports")
	println("  dump <local folder>     Dump the pyboard filesystem to the local folder")
	println("  sync <local folder>     Sync the local folder to the pyboard filesystem")
	println("  exec					   Execute python code on the pyboard")
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
