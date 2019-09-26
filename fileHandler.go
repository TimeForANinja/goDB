package main

import (
	"os"

	util "github.com/timeforaninja/goDB/utility"
)

type database struct {
	//	fileinfo      fileInfo
	file       *os.File
	passphrase string
	head       *head
}

func (db database) writeHead(h *head) error {
	c, err := h.serializeHead(util.StringtoBytes(db.passphrase))
	if err != nil {
		return err
	}
	db.file.WriteAt(c, 0)
	return nil
}

func (db database) readHead() (*head, error) {
	data := make([]byte, 128)
	db.file.ReadAt(data, 0)
	return deserializeHead(data, util.StringtoBytes(db.passphrase))
}

// NewDB is the factory for a new database
func NewDB() *database {
	return nil
}
func NewEncDB(userPW string, iv []byte) *database {
	userSecret := util.Hash(util.StringtoBytes(userPW), iv)
	userSecret[0] = userSecret[0]
	return nil
}

func main() {}
