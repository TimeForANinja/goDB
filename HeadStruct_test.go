package main

import (
	"testing"

	util "github.com/timeforaninja/goDB/utility"
)

func TestEncHead(t *testing.T) {
	key, _ := util.RandomIV(32)
	h := head{
		UseEncryption:          true,
		MasterPassword:         key,
		Version:                util.NewVersion(2, 3, 4),
		PageSize:               2048,
		PageCount:              420,
		TableListLocation:      69,
		EmptyPagesListLocation: 1337,
	}
	userPW := util.Hash([]byte{00, 01, 02, 03}, []byte{04, 05, 06, 07})
	// serialize head
	data, err := h.serializeHead(userPW)
	if err != nil {
		t.Error(err)
		return
	}
	// deserialize again
	h2, err := deserializeHead(data, userPW)
	if err != nil {
		t.Error(err)
		return
	}
	// compare
	if !h.equals(h2) {
		t.Error("objects not equal", h, h2)
	}
}

func TestHead(t *testing.T) {
	h := head{
		Version:                util.NewVersion(2, 3, 4),
		PageSize:               2048,
		PageCount:              420,
		TableListLocation:      69,
		EmptyPagesListLocation: 1337,
	}
	// serialize head
	data, err := h.serializeHead(nil)
	if err != nil {
		t.Error(err)
		return
	}
	// deserialize again
	h2, err := deserializeHead(data, nil)
	if err != nil {
		t.Error(err)
		return
	}
	// compare
	if !h.equals(h2) {
		t.Error("objects not equal", h, h2)
	}
}
