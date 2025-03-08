package tools

import (
	"bufio"
	"mputil/pyboard"
	"os"
)

func Tool_Exec(args []string, board *pyboard.Pyboard) {
	if board == nil {
		println("No device selected")
		return
	}
	// prompt the user for python code
	reader := bufio.NewReader(os.Stdin)
	println("Enter the code to eval:")
	print("> ")
	code, _ := reader.ReadString('\n')
	resp, err := board.Exec(code)
	if err {
		println("Error executing python code")
		println("----------------------------")
		println(resp)
	} else {
		println(resp)
	}
}
