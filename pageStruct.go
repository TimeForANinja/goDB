package goDB

import (
	"bytes"
	"errors"

	util "github.com/timeforaninja/goDB/utility"
)

// page represents the parsed pageheader + pagecontent
// page indexes start at 1, zero is the null pointer
type page struct {
	db       *Database
	index    uint32
	pageHead *pageHead
	data     []byte
}

func (page *page) equals(p2 *page) bool {
	return bytes.Equal(page.data, p2.data) &&
		page.index == p2.index &&
		page.pageHead.equals(p2.pageHead)
}

// readPage returns a new untrimmed page read from file by its index
func readPage(db *Database, pageNum uint32) (*page, error) {
	if pageNum == 0 {
		return nil, errors.New("nullpointer exception")
	}

	buffer, err := db.getPageBytes(pageNum)
	if err != nil {
		return nil, err
	}

	pageHead, data, err := deserializePage(buffer, db.useEncryption, db.dbHead.masterKey)
	if err != nil {
		return nil, err
	}
	return &page{db, pageNum, pageHead, data}, nil
}

// writePage writes a page to file
func (page *page) writePage(db *Database) error {
	data, err := page.serializePage(db.dbHead.pageSize, db.useEncryption, db.dbHead.masterKey)
	if err != nil {
		return err
	}
	return db.writePageBytes(page.index, data)
}

// deserializePage parses bytes of a page into a pageHead and the pageContent (handles encryption)
func deserializePage(data []byte, encrypted bool, masterKey []byte) (*pageHead, []byte, error) {
	headBytes := data[:48]
	ivBytes := data[48:64]
	dataBytes := data[64:]
	if !encrypted {
		return deserializePageHead(headBytes), dataBytes, nil
	}

	decHead, err := util.DecryptCFB(headBytes, ivBytes, masterKey, false)
	if err != nil {
		return nil, nil, err
	}
	decData, err := util.DecryptCFB(dataBytes, ivBytes, masterKey, false)
	if err != nil {
		return nil, nil, err
	}

	return deserializePageHead(decHead), decData, nil
}

// serializePage parses a page (pageHead and pageContent) into a single buffer (handles encryption)
func (page *page) serializePage(pageSize uint16, encryption bool, masterKey []byte) ([]byte, error) {
	data := make([]byte, pageSize)
	headBytes := page.pageHead.serializePageHead()

	if !encryption {
		copy(data[:48], headBytes)
		copy(data[64:], page.data)
		return data, nil
	}

	encHead, iv, err := util.NewEncryptCFB(headBytes, masterKey, false)
	if err != nil {
		return nil, err
	}
	encData, err := util.EncryptCFB(page.data, iv, masterKey, false)
	if err != nil {
		return nil, err
	}
	copy(data[:48], encHead)
	copy(data[48:64], iv)
	copy(data[64:], encData)
	return data, nil
}

// newPage creates a new page structure (from its head)
func newPage(db *Database, pageHead *pageHead) *page {
	p := page{
		db:       db,
		index:    db.dbHead.pageCount + 1,
		pageHead: pageHead,
		data:     make([]byte, db.dbHead.pageSize-64),
	}
	db.dbHead.pageCount++
	db.writeHead()
	return &p
}
