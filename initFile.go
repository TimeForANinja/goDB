package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
)

type Database struct {
	fileinfo      fileInfo
	file          *os.File
	encryption    bool
	encryptionKey []byte
}

func readPage(db *Database) *tablePage {
	var data []byte
	if db.encryption {
		data = readEncryptedPage(db)
	} else {
		data = readUnencryptedPage(db)
	}
	return parsePage(data)
}
func readEncryptedPage(db *Database) []byte {
	// read db.fileinfo.databasePageSize bytes
	// use last aes.BlockSize bytes as iv
	// return the decrypted data
	return nil
}
func readUnencryptedPage(db *Database) []byte {
	// read db.fileinfo.databasePageSize bytes
	return nil
}
func parsePage([]byte) *tablePage {
	return nil
}

type Version struct {
	Major int8
	Minor int8
	Patch int8
}

type fileInfo struct {
	Version              Version
	databasePageSize     int32
	fileChangeCounter    int64
	databasePages        int64
	defaultPageCacheSize int64
}

func newFileInfo(data []byte) fileInfo {
	// TODO: parse fileinfo in here
	return fileInfo{}
}

type tablePage struct {
	indexOfTable            int64
	numberOfColumns         int64
	indexOfColumnDefPage    int64
	indexOfFirstContentPage int64
	pageCount               int64
	rowCount                int64
	tableNameLength         int64
	tableName               string
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
	db.fileinfo = newFileInfo(data[5:])
	return db, nil
}
func OpenEnc(file, passphrase string) (*Database, error) {
	db := &Database{encryption: true}
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
	db.fileinfo = newFileInfo(dataDec[41:])
	return db, nil
}
func NewDBEnc(file, passphrase string) (*Database, error) {
	db := &Database{encryption: true}
	var err error

	db.file, err = os.Create(file)
	if err != nil {
		return nil, err
	}
	dataDec := make([]byte, 112)
	copy(dataDec[:9], []byte{103, 111, 68, 66, 32, 101, 110, 99, 00})
	db.encryptionKey, err = randomIV(32)
	if err != nil {
		return nil, err
	}
	copy(dataDec[9:41], db.encryptionKey)
	dataDec[41] = 69
	cipher, iv, err := newEncryptCFB(dataDec, passphrase)
	if err != nil {
		return nil, err
	}
	db.file.Write(append(iv[:], cipher[:]...))
	db.fileinfo = newFileInfo(dataDec[41:])
	return db, nil
}

func main() {
	_, err := NewDBEnc("./test.db.enc", "asdfMovie")
	fmt.Println("-- newDBEnc returned")
	fmt.Print("err: ")
	fmt.Println(err)
	_, err2 := OpenEnc("./test.db.enc", "asdfMovie")
	fmt.Println("-- OpenEnc returned")
	fmt.Print("err: ")
	fmt.Println(err2)
}

// cryptography methods
func hash(data, salt []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}
func randomIV(size int) ([]byte, error) {
	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}
func newEncryptCFB(plaintext []byte, passphrase string) ([]byte, []byte, error) {
	iv, err := randomIV(16)
	if err != nil {
		return nil, nil, err
	}
	cipher, err := encryptCFB(plaintext, iv, hash([]byte(passphrase), iv))
	if err != nil {
		return nil, nil, err
	}
	return cipher, iv, nil
}
func encryptCFB(plaintext, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext, nil
}
func decryptCFB(ciphertext, iv, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}
