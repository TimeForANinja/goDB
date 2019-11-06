package goDB

import (
	util "github.com/timeforaninja/goDB/utility"
)

type emptyTableListIterator struct {
	db       *database
	position *util.Position
	currPage *page
}

func (etli *emptyTableListIterator) forward() error {
	if (etli.position.Byte - etli.currPage.pageHead.endTrim) == etli.db.head.pageSize {
		var err error
		etli.position.Page = etli.currPage.pageHead.nextPage
		etli.currPage, err = ReadPage(etli.db, etli.currPage.pageHead.nextPage)
		if err != nil {
			return err
		}
		etli.position.Byte = etli.currPage.pageHead.firstItem
	} else {
		etli.position.Byte += 4
	}
	return nil
}

func (etli *emptyTableListIterator) read() uint32 {
	return util.BytesToUInt32(etli.currPage.data[etli.position.Byte:4])
}

func (etli *emptyTableListIterator) hasNext() bool {
	if (etli.position.Byte - etli.currPage.pageHead.endTrim) == etli.db.head.pageSize {
		return etli.currPage.pageHead.nextPage != 0
	}
	return true
}

func (db *database) getEmptyTableListIterator() (*emptyTableListIterator, error) {
	page, err := ReadPage(db, db.head.emptyPagesListLocation)
	if err != nil {
		return nil, err
	}
	pos := util.Position{
		Page: db.head.emptyPagesListLocation,
		Byte: page.pageHead.firstItem,
	}
	return &emptyTableListIterator{
		db:       db,
		position: &pos,
		currPage: page,
	}, nil
}
