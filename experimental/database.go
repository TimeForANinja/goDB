package main

import (
	"bytes"
	"errors"
	"os"
)

/*
 * main Database object & exported functions
 */
type Database struct {
	fileinfo      fileInfo
	file          *os.File
	passphrase    string
	encryption    bool
	encryptionKey []byte
}

func (db Database) writeHead() error {
	data := make([]byte, 128)
	if db.encryption {
		iv, err := randomIV(16)
		if err != nil {
			return err
		}
		copy(data[:16], iv)

		content := make([]byte, 112)
		copy(content[:9], []byte{103, 111, 68, 66, 32, 101, 110, 99, 00})
		copy(content[9:41], db.encryptionKey)
		copy(content[41:], db.fileinfo.toBytes())

		encryptedData, err := encryptCFB(content, iv, hash([]byte(db.passphrase), iv))
		if err != nil {
			return err
		}
		copy(data[16:], encryptedData)
	} else {
		copy(data[:5], []byte{103, 111, 68, 66, 00})
		copy(data[5:], db.fileinfo.toBytes())
	}

	db.file.WriteAt(data, 0)
	return nil
}

func Open(file string) (*Database, error) {
	db := &Database{encryption: false, encryptionKey: nil}
	var err error

	db.file, err = os.OpenFile(file, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	defer db.file.Close()

	data := make([]byte, 128)
	db.file.Read(data)
	if !bytes.Equal(data[:5], []byte{103, 111, 68, 66, 00}) {
		return nil, errors.New("Invalid File")
	}
	db.fileinfo = parseFileInfo(data[5:])
	return db, nil
}
func NewDB(file string) (*Database, error) {
	db := &Database{encryption: false, encryptionKey: nil}
	var err error

	db.file, err = os.Create(file)
	if err != nil {
		return nil, err
	}
	db.fileinfo = defaultFileInfo()
	err = db.writeHead()
	if err != nil {
		return nil, err
	}
	return db, nil
}
func OpenEnc(file, passphrase string) (*Database, error) {
	db := &Database{encryption: true, passphrase: passphrase}
	var err error

	db.file, err = os.OpenFile(file, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	defer db.file.Close()

	dataEnc := make([]byte, 128)
	db.file.Read(dataEnc)
	dataDec, err := decryptCFB(dataEnc[16:], dataEnc[:16], hash([]byte(passphrase), dataEnc[:16]))
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(dataDec[:9], []byte{103, 111, 68, 66, 32, 101, 110, 99, 00}) {
		return nil, errors.New("Invalid File")
	}
	db.encryptionKey = dataDec[9:41]
	db.fileinfo = parseFileInfo(dataDec[41:])
	return db, nil
}
func NewDBEnc(file, passphrase string) (*Database, error) {
	db := &Database{encryption: true, passphrase: passphrase}
	var err error

	db.file, err = os.Create(file)
	if err != nil {
		return nil, err
	}
	db.encryptionKey, err = randomIV(32)
	if err != nil {
		return nil, err
	}
	db.fileinfo = defaultFileInfo()
	err = db.writeHead()
	if err != nil {
		return nil, err
	}
	return db, nil
}
