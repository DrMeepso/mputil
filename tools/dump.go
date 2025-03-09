package tools

import (
	"fmt"
	"mputil/pyboard"
	"os"
)

// dump all files from the pyboard to a folder (src)
func Tool_Dump(args []string, board *pyboard.Pyboard, skipPrompt bool) {
	if len(args) != 2 {
		println("Usage: dump <folder>")
		return
	}

	// make sure the folder exists
	folder := args[1]
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		println("Error creating folder")
		return
	}

	// check if the folder is empty
	entries, err := os.ReadDir(folder)
	if err != nil {
		println("Error reading folder")
		return
	}
	if len(entries) > 0 && !skipPrompt {
		println("Folder is not empty, would you like to wipe it? (y/n)")
		print("> ")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			println("Error reading response")
		}
		if response != "y" {
			println("Aborted")
			return
		}

		// remove all files in the folder
		for _, entry := range entries {
			err := os.Remove(folder + "/" + entry.Name())
			if err != nil {
				println("Error removing file", entry.Name())
			}
		}
	}

	println("Dumping files from the pyboard to", folder)

	// get the list of files
	files := board.FS.ListDir()
	longestName := 0
	for _, file := range files {
		if len(file) > longestName {
			longestName = len(file)
		}
	}

	padto := max(20, longestName+2)

	for i, file := range files {
		println("Dumping", padString(file, padto), fmt.Sprint(i+1)+"/"+fmt.Sprint(len(files)))
		content := board.FS.ReadFile(file)
		filePath := folder + "/" + file
		err := os.WriteFile(filePath, []byte(content), os.ModePerm)
		if err != nil {
			println("Error writing file", file)
		}
	}

	println("Dumped", len(files), "files to", folder)
}

func padString(input string, length int) string {
	if len(input) >= length {
		return input
	}
	padding := length - len(input)
	return input + fmt.Sprintf("%*s", padding, "")
}
