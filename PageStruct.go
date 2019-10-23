package goDB

import (
	"bytes"

	util "github.com/timeforaninja/goDB/utility"
)

// TypeTableList holds the id of a table_List page
const TypeTableList = 0

// TypeTableSchema holds the id of a table_schema page
const TypeTableSchema = 1

// TypeEmptyPagesList holds the id of a empty_pages_list page
const TypeEmptyPagesList = 2

// TypeTableRows holds the id of a table_rows page
const TypeTableRows = 3

type page struct {
	index    uint32
	pageHead *pageHead
	data     []byte
}

func (p *page) equals(p2 *page) bool {
	return bytes.Equal(p.data, p2.data) &&
		p.index == p2.index &&
		p.pageHead.equals(p2.pageHead)
}

type pageHead struct {
	pageType     uint8
	nextPage     uint32
	prevPage     uint32
	firstItem    uint16
	endTrim      uint16
	pageSpezific []byte
}

func (p *pageHead) equals(p2 *pageHead) bool {
	return p.pageType == p2.pageType &&
		p.nextPage == p2.nextPage &&
		p.prevPage == p2.prevPage &&
		p.firstItem == p2.firstItem &&
		p.endTrim == p2.endTrim &&
		bytes.Equal(p.pageSpezific, p2.pageSpezific)
}

func ReadPage(db *database, pageNum uint32) (*page, error) {
	buffer, err := getPageBytes(db, pageNum)
	if err != nil {
		return nil, err
	}

	pageHead, data, err := deserializePage(buffer, db.useEncryption, db.head.masterKey)
	if err != nil {
		return nil, err
	}
	return &page{pageNum, pageHead, data}, nil
}

func (page *page) WritePage(db *database) error {
	data, err := page.serializePage(db.head.pageSize, db.useEncryption, db.head.masterKey)
	if err != nil {
		return err
	}
	pageStart := 128 + ((int64(page.index) - 1) * int64(db.head.pageSize))
	_, err = db.file.WriteAt(data, pageStart)
	return err
}

func deserializePage(data []byte, encrypted bool, masterKey []byte) (*pageHead, []byte, error) {
	if !encrypted {
		return deserializePageHeadCore(data[:48]), data[64:], nil
	}

	iv := data[48:64]
	decHead, err := util.DecryptCFB(data[:48], iv, masterKey, false)
	if err != nil {
		return nil, nil, err
	}
	decData, err := util.DecryptCFB(data[64:], iv, masterKey, false)
	if err != nil {
		return nil, nil, err
	}

	return deserializePageHeadCore(decHead), decData, nil
}

func (page *page) serializePage(pageSize uint16, encryption bool, masterKey []byte) ([]byte, error) {
	data := make([]byte, pageSize)
	headCore := page.pageHead.serializePageHeadCore()

	if !encryption {
		copy(data[:48], headCore)
		copy(data[64:], page.data)
		return data, nil
	}

	encHead, iv, err := util.NewEncryptCFB(headCore, masterKey, false)
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

func deserializePageHeadCore(data []byte) *pageHead {
	head := pageHead{}
	head.pageType = uint8(data[0])
	head.nextPage = util.BytesToUInt32(data[1:5])
	head.prevPage = util.BytesToUInt32(data[5:9])
	head.firstItem = util.BytesToUInt16(data[9:11])
	head.endTrim = util.BytesToUInt16(data[11:13])
	head.pageSpezific = data[13:48]
	return &head
}

func (head *pageHead) serializePageHeadCore() []byte {
	data := make([]byte, 64)
	data[0] = head.pageType
	copy(data[1:5], util.Uint32toBytes(head.nextPage))
	copy(data[5:9], util.Uint32toBytes(head.prevPage))
	copy(data[9:11], util.Uint16toBytes(head.firstItem))
	copy(data[11:13], util.Uint16toBytes(head.endTrim))
	copy(data[13:48], head.pageSpezific)
	return data
}

func getPageBytes(db *database, pageNum uint32) ([]byte, error) {
	pageBuffer := make([]byte, db.head.pageSize)
	pageStart := 128 + ((int64(pageNum) - 1) * int64(db.head.pageSize))
	_, err := db.file.ReadAt(pageBuffer, pageStart)
	if err != nil {
		return nil, err
	}
	return pageBuffer, nil
}

func NewPage(db *database) *page {
	p := page{index: db.head.pageCount + 1, data: make([]byte, db.head.pageSize-64)}
	db.head.pageCount++
	db.writeHead(db.head)
	return &p
}
