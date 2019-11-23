package goDB

import (
	"os"

	util "github.com/timeforaninja/goDB/utility"
)

// Database is the general struct you interact with
type Database struct {
	file          *os.File
	useEncryption bool
	secret        *secret
	dbHead        *dbHead
}

type secret struct {
	userPassphrase []byte // the string
}

// writeHead writes the database head to the file
func (db *Database) writeHead() error {
	c, err := db.dbHead.serializeHead(db.secret.userPassphrase)
	if err != nil {
		return err
	}
	db.file.WriteAt(c, 0)
	return nil
}

// readHead reads a database head from file
func (db *Database) readHead() (*dbHead, error) {
	data := make([]byte, 128)
	db.file.ReadAt(data, 0)
	return deserializeHead(data, db.secret.userPassphrase)
}

// getPageBytes reads the pagehead and page bytes from file (doesnt handle encryption)
func (db *Database) getPageBytes(pageNum uint32) ([]byte, error) {
	pageBuffer := make([]byte, db.dbHead.pageSize)
	pageStart := 128 + ((int64(pageNum) - 1) * int64(db.dbHead.pageSize))
	_, err := db.file.ReadAt(pageBuffer, pageStart)
	if err != nil {
		return nil, err
	}
	return pageBuffer, nil
}

// writePageBytes writes pagehead and page bytes to file (doesnt handle encryption)
func (db *Database) writePageBytes(pageNum uint32, data []byte) error {
	pageStart := 128 + ((int64(pageNum) - 1) * int64(db.dbHead.pageSize))
	_, err := db.file.WriteAt(data, pageStart)
	return err
}

// Open initialises a database from the provided file
func (db *Database) Open(filename string) error {
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
	db.dbHead, err = db.readHead()
	return err
}

// newFile is used when calling Open with a file that doesn't yet exist
// newFile also initialises a new blank header
func (db *Database) newFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	db.file = file
	db.dbHead, err = newBlankHead(db.useEncryption)
	if err != nil {
		return err
	}
	err = db.writeHead()
	return err
}

// Vacuum optimizes the database by removing mostly trim spaces and empty pages
func (db *Database) Vacuum() {
	// TODO: clean up the database
}

// NewDB is the factory for a new database
func NewDB() *Database {
	db := Database{useEncryption: false, secret: nil}
	return &db
}

// NewEncDB is the factory for a new encrypted database
func NewEncDB(passphrase string) *Database {
	db := Database{
		useEncryption: true,
		secret:        &secret{util.StringtoBytes(passphrase)},
	}
	return &db
}
