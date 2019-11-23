package goDB

import (
	util "github.com/timeforaninja/goDB/utility"
)

type pageIterator struct {
	db       *Database
	position *util.Position
	currPage *page
}

func (pi *pageIterator) top(length uint8) []byte {
	return nil
}

func (pi *pageIterator) pop(length uint8) []byte {
	return nil
}

// New creates a new pageIterator
func (db *Database) newPageIterator(pageNum uint32) (*pageIterator, error) {
	page, err := readPage(db, pageNum)
	if err != nil {
		return nil, err
	}
	return &pageIterator{
		db:       db,
		currPage: page,
	}, nil
}
