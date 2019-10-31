package goDB

import (
	util "github.com/timeforaninja/goDB/utility"
)

type emptyTableListIterator struct {
	db       *database
	position *util.Position
	currPage *page
}

func (etli *emptyTableListIterator) next() uint32 {
	if(etli.position.Byte == etli.db.head.pageSize) // do sth special
	etli.position.Byte += 4
	return util.BytesToUInt32(etli.currPage.data[etli.position.Byte-4 : 4])
}

func (db *database) getEmptyTableListIterator() *emptyTableListIterator {
	return emptyTableListIterator{
		db:       db,
		position: Position{db.head.emptyPagesListLocation},
		currPage: ReadPage(db, db.head.emptyPagesListLocation),
	}
}
