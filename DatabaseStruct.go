package main

import (
	"os"

	util "github.com/timeforaninja/goDB/utility"
)

type database struct {
	//	fileinfo      fileInfo
	file           *os.File
	useEncryption  bool
	userPassphrase []byte // the string
}

func (db database) writeHead(h *head) error {
	c, err := h.serializeHead(db.userPassphrase)
	if err != nil {
		return err
	}
	db.file.WriteAt(c, 0)
	return nil
}

func (db database) readHead() (*head, error) {
	data := make([]byte, 128)
	db.file.ReadAt(data, 0)
	return deserializeHead(data, db.userPassphrase)
}

// NewDB is the factory for a new database
func NewDB() *database {
	db := database{useEncryption: false}
	return &db
}

// NewEncDB is the factory for a new encrypted database
func NewEncDB(passphrase string) *database {
	db := database{useEncryption: true, userPassphrase: util.StringtoBytes(passphrase)}
	return &db
}

func main() {}
