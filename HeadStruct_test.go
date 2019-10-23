package goDB

import (
	"testing"

	util "github.com/timeforaninja/goDB/utility"
)

func TestEncHead(t *testing.T) {
	testKey, _ := util.RandomIV(32)
	testHead := head{
		useEncryption:          true,
		masterKey:              testKey,
		version:                util.NewVersion(2, 3, 4),
		pageSize:               2048,
		pageCount:              420,
		tableListLocation:      69,
		emptyPagesListLocation: 1337,
	}
	testUserPassphrase := []byte{00, 01, 02, 03}

	// serialize head
	testData, err := testHead.serializeHead(testUserPassphrase)
	if err != nil {
		t.Error(err)
		return
	}
	if len(testData) != 128 {
		t.Error("testData invalid length", testData, 128)
	}
	// deserialize again
	testHeadDecoded, err := deserializeHead(testData, testUserPassphrase)
	if err != nil {
		t.Error(err)
		return
	}
	// compare
	if !testHead.equals(testHeadDecoded) {
		t.Error("objects not equal", testHead, testHeadDecoded)
	}
}

func TestHead(t *testing.T) {
	testHead := head{
		version:                util.NewVersion(2, 3, 4),
		pageSize:               2048,
		pageCount:              420,
		tableListLocation:      69,
		emptyPagesListLocation: 1337,
	}
	// serialize head
	testData, err := testHead.serializeHead(nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(testData) != 128 {
		t.Error("testData invalid length", testData, 128)
	}
	// deserialize again
	testHeadDecoded, err := deserializeHead(testData, nil)
	if err != nil {
		t.Error(err)
		return
	}
	// compare
	if !testHead.equals(testHeadDecoded) {
		t.Error("objects not equal", testHead, testHeadDecoded)
	}
}
