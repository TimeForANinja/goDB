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

// pageHead represents the header (information) of a page
type pageHead struct {
	pageType     uint8
	nextPage     uint32
	prevPage     uint32
	firstItem    uint16
	endTrim      uint16
	pageSpezific []byte
}

func (head *pageHead) equals(p2 *pageHead) bool {
	return head.pageType == p2.pageType &&
		head.nextPage == p2.nextPage &&
		head.prevPage == p2.prevPage &&
		head.firstItem == p2.firstItem &&
		head.endTrim == p2.endTrim &&
		bytes.Equal(head.pageSpezific, p2.pageSpezific)
}

// deserializePageHead parses bytes into a pageHead object
func deserializePageHead(data []byte) *pageHead {
	head := pageHead{}
	head.pageType = uint8(data[0])
	head.nextPage = util.BytesToUInt32(data[1:5])
	head.prevPage = util.BytesToUInt32(data[5:9])
	head.firstItem = util.BytesToUInt16(data[9:11])
	head.endTrim = util.BytesToUInt16(data[11:13])
	head.pageSpezific = data[13:48]
	return &head
}

// serializePageHead parses a pageHead into its bytes
func (head *pageHead) serializePageHead() []byte {
	data := make([]byte, 64) // might change to 48 since the iv gets placed into a new buffer anyway
	data[0] = head.pageType
	copy(data[1:5], util.Uint32toBytes(head.nextPage))
	copy(data[5:9], util.Uint32toBytes(head.prevPage))
	copy(data[9:11], util.Uint16toBytes(head.firstItem))
	copy(data[11:13], util.Uint16toBytes(head.endTrim))
	copy(data[13:48], head.pageSpezific)
	return data
}
