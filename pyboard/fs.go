package pyboard

import (
	"fmt"
	"strings"
)

type PyFileSystem struct {
	pyboard    *Pyboard
	CurrentDir string
}

func NewPyFileSystem() *PyFileSystem {
	return &PyFileSystem{
		CurrentDir: "/",
	}
}

func (fs *PyFileSystem) ListDir() []string {
	files, _ := fs.pyboard.Exec("import os; print(','.join(os.listdir()))")
	return strings.Split(files, ",")
}

func (fs *PyFileSystem) ChangeDir(dir string) {
	fs.pyboard.Exec("import os; os.chdir('" + dir + "')")
	fs.CurrentDir = dir
}

// read the file in mutiple chunks
func (fs *PyFileSystem) readFileChunked(filename string, chunkSize int) string {
	python := "import os\n\r"
	python += "fileData = b''\n\r"
	python += "with open('" + filename + "', 'rb') as f:\n\r"
	python += "    while True:\n\r"
	python += "        data = f.read(" + fmt.Sprint(chunkSize) + ")\n\r"
	python += "        if not data:\n\r"
	python += "            break\n\r"
	python += "        fileData += data\n\r"
	python += "print(str(fileData, 'utf-8'))\n\r"

	fileContent, _ := fs.pyboard.Exec(python)

	return fileContent
}

func (fs *PyFileSystem) ReadFile(filename string) string {
	return fs.readFileChunked(filename, 256)
}
