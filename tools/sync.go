package tools

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mputil/pyboard"
	"os"
)

const (
	FileWaiting    = 0 // for a file we are yet to check
	FileHashChange = 1 // for a file that has changed
	FileMissing    = 2 // for a file that is missing in the local folder
	FileExtra      = 3 // for a file that is extra in the local folder
	FileSynced     = 4 // for a file that is synced
)

const (
	CONSOLE_Reset  = "\033[0m"
	CONSOLE_Red    = "\033[31m"
	CONSOLE_Green  = "\033[32m"
	CONSOLE_Yellow = "\033[33m"
	CONSOLE_Blue   = "\033[34m"
)

func Tool_Sync(args []string, board *pyboard.Pyboard) {
	if len(args) < 2 {
		println("Usage: sync <src>")
		return
	}

	src := args[1]

	// check if the folder exists
	_, err := os.Stat(src)
	if err != nil {
		println("Folder does not exist")
		return
	}

	allFiles := make(map[string]int)

	// get the list of files
	localFiles, err := os.ReadDir(src)
	if err != nil {
		println("Error reading folder")
		return
	}

	for _, file := range localFiles {
		allFiles[file.Name()] = FileExtra
	}

	println("Checking files...")

	// get the list of files
	files := board.FS.ListDir()
	for _, file := range files {
		// check if the file exists in the src folder
		_, err := os.Stat(src + "/" + file)
		if err != nil {
			allFiles[file] = FileMissing
			continue
		}

		// get a sha256 hash of the file
		hash := board.FS.GetSHA256(file)
		if hash == "" {
			continue
		}

		// get a sha256 hash of the file in the src folder
		hashSrc, err := getFileSHA256(src + "/" + file)
		if err != nil {
			continue
		}

		// compare the hashes
		if hash != hashSrc {
			allFiles[file] = FileHashChange
			continue
		}
		allFiles[file] = FileSynced
	}

	// print the results
	unsureFiles := int(0)
	wrongFiles := (make(map[string]int))
	for file, status := range allFiles {
		switch status {
		case FileHashChange:
			print(CONSOLE_Yellow)
			println("Changed >", file)
			print(CONSOLE_Reset)
			unsureFiles++
			wrongFiles[file] = FileHashChange
		case FileMissing:
			print(CONSOLE_Red)
			println("Missing >", file)
			print(CONSOLE_Reset)
			unsureFiles++
			wrongFiles[file] = FileMissing
		case FileExtra:
			print(CONSOLE_Blue)
			println("Extra >", file)
			print(CONSOLE_Reset)
			unsureFiles++
			wrongFiles[file] = FileExtra
		}
	}

	if unsureFiles > 0 {
		println("^-- Out of sync files --^")
		println("You have some file out of sync with the pyboard")
		println("1 > Sync files from the pyboard")
		println("2 > Upload changed, missing or extra files?")
		println("3 > Cancel")
		print("> ")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			println("Error reading response")
		}

		switch response {
		case "1":
			// sync files from the pyboard
			Tool_Dump([]string{"dump", src}, board, true)
			println("Files synced")
		case "2":
			// upload changed, missing or extra files
			for fileName, file := range wrongFiles {
				switch file {
				case FileHashChange:
					println("Uploading changed file", fileName)
					localContent, err := os.ReadFile(src + "/" + fileName)
					if err != nil {
						println("Error reading file", fileName)
						continue
					}
					board.FS.WriteFile(fileName, string(localContent))
				case FileMissing:
					// delete the file from the pyboard
					println("Deleting missing file", fileName)
					board.FS.RemoveFile(fileName)
				case FileExtra:
					println("Uploading extra file", fileName)
					localContent, err := os.ReadFile(src + "/" + fileName)
					if err != nil {
						println("Error reading file", fileName)
						continue
					}
					board.FS.WriteFile(fileName, string(localContent))
				}
			}
			println("Files uploaded")

		case "3":
			println("Aborted")
			return
		}
	}

	println("Sync has begun")

}

func getFileSHA256(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
