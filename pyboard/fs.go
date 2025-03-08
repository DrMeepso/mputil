package pyboard

import (
	"fmt"
	"strconv"
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

// write the file in mutiple chunks
func (fs *PyFileSystem) writeFileChunked(filename string, data string, chunkSize int) bool {
	fileData := []byte(data)
	chunkCount := len(fileData) / chunkSize

	for i := 0; i < chunkCount; i++ {
		chunk := fileData[i*chunkSize : (i+1)*chunkSize]
		python := "import os\n\r"
		python += "with open('" + filename + "', 'wb') as f:\n\r"
		python += "    f.write(" + fmt.Sprintf("%q", chunk) + ")\n\r"
		ret, err := fs.pyboard.Exec(python)
		println(ret)
		if err {
			return false
		}
	}

	// write the last chunk
	if len(fileData)%chunkSize != 0 {
		chunk := fileData[chunkCount*chunkSize:]
		python := "import os\n\r"
		python += "with open('" + filename + "', 'wb') as f:\n\r"
		python += "    f.write(" + fmt.Sprintf("%q", chunk) + ")\n\r"
		fs.pyboard.Exec(python)
		ret, err := fs.pyboard.Exec(python)
		println(ret)
		if err {
			return false
		}
	}

	return true
}

func (fs *PyFileSystem) WriteFile(filename string, data string) bool {
	return fs.writeFileChunked(filename, data, 256)
}

func (fs *PyFileSystem) RemoveFile(filename string) {
	fs.pyboard.Exec("import os; os.remove('" + filename + "')")
}

func (fs *PyFileSystem) MakeDir(dir string) {
	fs.pyboard.Exec("import os; os.mkdir('" + dir + "')")
}

func (fs *PyFileSystem) RemoveDir(dir string) {
	fs.pyboard.Exec("import os; os.rmdir('" + dir + "')")
}

func (fs *PyFileSystem) GetSize(filename string) int {
	size, _ := fs.pyboard.Exec("import os; print(os.stat('" + filename + "').st_size)")
	fileSize, err := strconv.Atoi(size)
	if err != nil {
		return -1
	}
	return fileSize
}

func (fs *PyFileSystem) Exists(filename string) bool {
	ret, _ := fs.pyboard.Exec("import os; print(os.path.exists('" + filename + "'))")
	return ret == "True"
}

func (fs *PyFileSystem) IsDir(filename string) bool {
	ret, _ := fs.pyboard.Exec("import os; print(os.path.isdir('" + filename + "'))")
	return ret == "True"
}

func (fs *PyFileSystem) GetSHA256(filename string) string {

	python := "import hashlib\n\r"
	python += "import binascii\n\r"
	python += "hash = hashlib.sha256()\n\r"
	python += "with open('" + filename + "', 'rb') as f:\n\r"
	python += "    while True:\n\r"
	python += "        data = f.read(256)\n\r"
	python += "        if not data:\n\r"
	python += "            break\n\r"
	python += "        hash.update(data)\n\r"
	python += "print(binascii.hexlify(hash.digest()).decode('utf-8'))\n\r"

	hash, err := fs.pyboard.Exec(python)
	if err {
		println(hash)
		return ""
	}
	return hash
}
