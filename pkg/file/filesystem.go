package file

import (
	"os"
	"encoding/json"
	"bytes"
)

func CreateDir(path string) error {
	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		if err := os.MkdirAll(path, 0777); err != nil {
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}

func WriteJSON(fileName string, data []byte, outDir string) error {
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.Mkdir(outDir, 0777)
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
