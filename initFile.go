package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"os"
)

type version struct {
	major int8
	minor int8
	patch int8
}

/*
 * everything about a page in the database
 */
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

/*
 * data transformation utilities
 */
func int16toBytes(int int16) []byte {
	return []byte{
		byte(int >> 8),
		byte(int & 255),
	}
}
func int32toBytes(int int32) []byte {
	return []byte{
		byte((int >> 24) & 255),
		byte((int >> 16) & 255),
		byte((int >> 8) & 255),
		byte(int & 255),
	}
}
func bytesToInt16(bytes []byte) int16 {
	return (int16(bytes[0]) << 8) + int16(bytes[1])
}
func bytesToInt32(bytes []byte) int32 {
	return (int32(bytes[2]) << 24) + (int32(bytes[2]) << 16) + (int32(bytes[2]) << 8) + int32(bytes[3])
}

/*
 * utility for the meta information in the first 128 bytes
 */
type fileInfo struct {
	version              version
	databasePageSize     int16
	fileChangeCounter    int32
	databasePages        int32
	defaultPageCacheSize int32
}

func (f fileInfo) toBytes() []byte {
	resp := make([]byte, 17)
	resp[0] = byte(f.version.major)
	resp[1] = byte(f.version.minor)
	resp[2] = byte(f.version.patch)
	copy(resp[3:5], int16toBytes(f.databasePageSize))
	copy(resp[5:9], int32toBytes(f.fileChangeCounter))
	copy(resp[9:13], int32toBytes(f.databasePages))
	copy(resp[13:17], int32toBytes(f.defaultPageCacheSize))
	return resp
}

func parseFileInfo(data []byte) fileInfo {
	f := fileInfo{}
	f.version = version{
		major: int8(data[0]),
		minor: int8(data[1]),
		patch: int8(data[2]),
	}
	f.databasePageSize = bytesToInt16(data[3:5])
	f.fileChangeCounter = bytesToInt32(data[5:9])
	f.databasePages = bytesToInt32(data[9:13])
	f.defaultPageCacheSize = bytesToInt32(data[13:17])
	return f
}
func defaultFileInfo() fileInfo {
	return fileInfo{
		version:              version{major: 0, minor: 0, patch: 1},
		databasePageSize:     4096,
		fileChangeCounter:    0,
		databasePages:        0,
		defaultPageCacheSize: 0,
	}
}

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

/*
 * bullshit functions
 */
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

/*
 * cryptography functions
 */
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
