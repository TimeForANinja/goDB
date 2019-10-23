package goDB

import (
	util "github.com/timeforaninja/goDB/utility"
)

type emptyTableList struct {
	page        *page
	emptyTables []uint32
}

func (page *page) parseAsEmptyTableList() *emptyTableList {
	tables := make([]uint32, 0, len(page.data))
	dataLength := (len(page.data) - int(page.pageHead.endTrim)) / 4
	for i := 0; i < dataLength; i++ {
		tables[i] = util.BytesToUInt32(page.data[i*4 : (i+1)*4])
	}
	return &emptyTableList{page: page, emptyTables: tables}
}
