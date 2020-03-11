package file

import "os"

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
