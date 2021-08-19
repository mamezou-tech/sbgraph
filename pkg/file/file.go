package file

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

// CreateDir will create directory according to path argument
func CreateDir(path string) error {
	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		if err := os.MkdirAll(path, 0775); err != nil {
			return err
		}
	}
	return nil
}

// WriteBytes will write bytes to file
func WriteBytes(data []byte, fileName string, outDir string) error {
	dir, err := os.Stat(outDir)
	if os.IsNotExist(err) || !dir.IsDir() {
		return err
	}
	file, err := os.Create(outDir + "/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	var pj bytes.Buffer
	json.Indent(&pj, []byte(data), "", " ")
	file.Write(pj.Bytes())
	return nil
}

// ReadBytes will read bytes from file
func ReadBytes(fileName string, outDir string) ([]byte, error) {
	raw, err := ioutil.ReadFile(outDir + "/" + fileName)
	if err != nil {
		return nil, err
	}
	return raw, err
}

// Exists will check if the file exists
func Exists(fileName string, outDir string) bool {
	_, err := os.Stat(outDir + "/" + fileName)
	return err == nil
}
