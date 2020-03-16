package file

import (
	"os"
	"encoding/json"
	"bytes"
)

func CreateDir(path string) error {
	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		if err := os.MkdirAll(path, 0775); err != nil {
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}

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
