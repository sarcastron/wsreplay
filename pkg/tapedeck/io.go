package tapedeck

import (
	"encoding/gob"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// Tape files are just gob files but this analogy ain't tired yet
func WriteTape(filePath string, messages *[]Message) error {
	baseDir := filepath.Dir(filePath)
	exists, err := exists(baseDir)
	if err != nil {
		return err
	}
	if !exists {
		err = os.Mkdir(baseDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(messages)
	file.Close()
	return nil
}

func ReadTape(filePath string, messages *[]Message) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(messages)
	}
	file.Close()
	return err
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
