package goDB

import (
	util "github.com/timeforaninja/goDB/utility"
)

type table struct {
	uid                         uint32
	pageindexOfFirstTableSchema uint32
	pageindexOfFirstTableRows   uint32
	pageindexOfLastTableRows    uint32
	rowCount                    uint32
	columnCount                 uint32
	name                        string
}

type tableList struct {
	page         *page
	tableCount   uint32
	tables       []*table
	overFlowData []byte
}

func (page *page) parseAsTableList() *tableList {
	tc := util.BytesToUInt32(page.pageHead.pageSpezific[:4])
	var overflow []byte
	if page.pageHead.firstItem == 0 {
		overflow = page.data
		return &tableList{page: page, tableCount: tc, overFlowData: overflow}
	}
	tables := make([]*table, 0, tc)
	pos := page.pageHead.firstItem - 1 // switch to 0 indexing
	for i := 0; true; i++ {
		length := util.BytesToUInt32(page.data[pos : pos+4])
		if uint(pos)+uint(length) > uint(len(page.data))-uint(page.pageHead.endTrim) {
			overflow = page.data[pos:]
			break
		}
		t := table{}
		t.uid = util.BytesToUInt32(page.data[pos+4 : pos+8])
		t.pageindexOfFirstTableSchema = util.BytesToUInt32(page.data[pos+8 : pos+12])
		t.pageindexOfFirstTableRows = util.BytesToUInt32(page.data[pos+12 : pos+16])
		t.pageindexOfLastTableRows = util.BytesToUInt32(page.data[pos+16 : pos+20])
		t.rowCount = util.BytesToUInt32(page.data[pos+20 : pos+24])
		t.columnCount = util.BytesToUInt32(page.data[pos+24 : pos+28])
		nameLength := util.BytesToUInt16(page.data[pos+28 : pos+30])
		t.name = util.BytesToString(page.data[pos+30 : pos+30+nameLength])
		tables[i] = &t
		pos += 30 + nameLength
	}
	return &tableList{page: page, tableCount: tc, overFlowData: overflow, tables: tables}
}
