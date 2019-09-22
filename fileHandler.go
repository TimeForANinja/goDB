package main

import (
	"os"
)

type Database struct {
	//	fileinfo      fileInfo
	file       *os.File
	passphrase string
	head       *Head
}

func (db Database) writeHead(h *Head) error {
	c, err := h.SerializeHead(stringtoBytes(db.passphrase))
	if err != nil {
		return err
	}
	db.file.WriteAt(c, 0)
	return nil
}

func (db Database) readHead() (*Head, error) {
	data := make([]byte, 128)
	db.file.ReadAt(data, 0)
	return DeserializeHead(data, stringtoBytes(db.passphrase))
}

func main() {}
