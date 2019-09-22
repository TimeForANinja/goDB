package main

import (
	"testing"
)

func TestEncHead(t *testing.T) {
	key, _ := randomIV(32)
	h := Head{
		UseEncryption:          true,
		MasterPassword:         key,
		Version:                Version{2, 3, 4},
		PageSize:               2048,
		PageCount:              420,
		TableListLocation:      69,
		EmptyPagesListLocation: 1337,
	}
	userPW := hash([]byte{00, 01, 02, 03}, []byte{04, 05, 06, 07})
	// serialize head
	data, err := h.SerializeHead(userPW)
	if err != nil {
		t.Error(err)
		return
	}
	// deserialize again
	h2, err := DeserializeHead(data, userPW)
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
	h := Head{
		Version:                Version{2, 3, 4},
		PageSize:               2048,
		PageCount:              420,
		TableListLocation:      69,
		EmptyPagesListLocation: 1337,
	}
	// serialize head
	data, err := h.SerializeHead(nil)
	if err != nil {
		t.Error(err)
		return
	}
	// deserialize again
	h2, err := DeserializeHead(data, nil)
	if err != nil {
		t.Error(err)
		return
	}
	// compare
	if !h.equals(h2) {
		t.Error("objects not equal", h, h2)
	}
}
