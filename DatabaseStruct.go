package goDB

import (
	"os"

	util "github.com/timeforaninja/goDB/utility"
)

type database struct {
	file           *os.File
	useEncryption  bool
	userPassphrase []byte // the string
	head           *head
}

func (db *database) writeHead(h *head) error {
	c, err := h.serializeHead(db.userPassphrase)
	if err != nil {
		return err
	}
	db.file.WriteAt(c, 0)
	return nil
}

func (db *database) readHead() (*head, error) {
	data := make([]byte, 128)
	db.file.ReadAt(data, 0)
	return deserializeHead(data, db.userPassphrase)
}

func (db *database) Open(filename string) error {
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return db.newFile(filename)
	}
	if err != nil {
		return err
	}

	if !stat.Mode().IsRegular() {
		return db.newFile(filename)
	}
	db.file, err = os.Open(filename)
	if err != nil {
		return err
	}
	db.head, err = db.readHead()
	return err
}

func (db *database) newFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	db.file = file
	db.head, err = newBlankHead(db.useEncryption)
	if err != nil {
		return err
	}
	err = db.writeHead(db.head)
	return err
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
