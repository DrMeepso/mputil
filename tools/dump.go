package tools

import (
	"fmt"
	"mputil/pyboard"
	"os"
)

// dump all files from the pyboard to a folder (src)
func Tool_Dump(args []string, board *pyboard.Pyboard) {
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
	if len(entries) > 0 {
		println("Folder is not empty, would you like to wipe it? (y/n)")
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

	// get the list of files
	files := board.FS.ListDir()
	for _, file := range files {
		println("Dumping", file)
		content := board.FS.ReadFile(file)
		filePath := folder + "/" + file
		err := os.WriteFile(filePath, []byte(content), os.ModePerm)
		if err != nil {
			println("Error writing file", file)
		}
	}

	println("Done")

}
