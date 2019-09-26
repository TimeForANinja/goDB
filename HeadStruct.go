package goDB

import (
	"bytes"
	"errors"

	util "github.com/timeforaninja/goDB/utility"
)

type head struct {
	useEncryption          bool
	masterKey              []byte
	version                util.Version
	pageSize               uint16
	pageCount              uint32
	tableListLocation      uint32
	emptyPagesListLocation uint32
}

func (h *head) equals(h2 *head) bool {
	return h.useEncryption == h2.useEncryption &&
		bytes.Equal(h.masterKey, h2.masterKey) &&
		h.version.Equals(h2.version) &&
		h.pageSize == h2.pageSize &&
		h.pageCount == h2.pageCount &&
		h.tableListLocation == h2.tableListLocation &&
		h.emptyPagesListLocation == h2.emptyPagesListLocation
}

func (h *head) serializeHeadCore() []byte {
	resp := make([]byte, 17)
	copy(resp[:3], h.version.ToBytes())
	copy(resp[3:5], util.Uint16toBytes(h.pageSize))
	copy(resp[5:9], util.Uint32toBytes(h.pageCount))
	copy(resp[9:13], util.Uint32toBytes(h.tableListLocation))
	copy(resp[13:17], util.Uint32toBytes(h.emptyPagesListLocation))
	return resp
}

func deserializeHeadCore(data []byte) *head {
	h := head{}
	h.version = util.NewVersionFromBytes(data[0:3])
	h.pageSize = util.BytesToUInt16(data[3:5])
	h.pageCount = util.BytesToUInt32(data[5:9])
	h.tableListLocation = util.BytesToUInt32(data[9:13])
	h.emptyPagesListLocation = util.BytesToUInt32(data[13:17])
	return &h
}

func (h *head) serializeHead(userPW []byte) ([]byte, error) {
	serializedHead := h.serializeHeadCore()
	fileHeader := make([]byte, 128)
	if !h.useEncryption {
		copy(fileHeader[:5], []byte{103, 111, 68, 66, 00})
		copy(fileHeader[5:], serializedHead)
		return fileHeader, nil
	}

	copy(fileHeader[:9], []byte{103, 111, 68, 66, 32, 101, 110, 99, 00})
	encMasterKey, iv, err := util.NewEncryptCFB(h.masterKey, userPW, true)
	if err != nil {
		return nil, err
	}
	copy(fileHeader[9:25], iv)
	copy(fileHeader[25:57], encMasterKey)
	testString, err := util.EncryptCFB([]byte{116, 114, 117, 101}, iv, h.masterKey, false)
	if err != nil {
		return nil, err
	}
	copy(fileHeader[57:61], testString)
	encData, err := util.EncryptCFB(serializedHead, iv, h.masterKey, false)
	if err != nil {
		return nil, err
	}
	copy(fileHeader[61:], encData)
	return fileHeader, nil
}

func deserializeHead(data []byte, userPW []byte) (*head, error) {
	if bytes.Equal(data[:5], []byte{103, 111, 68, 66, 00}) {
		head := deserializeHeadCore(data[5:])
		head.useEncryption = false
		return head, nil
	}

	if !bytes.Equal(data[:9], []byte{103, 111, 68, 66, 32, 101, 110, 99, 00}) {
		return nil, errors.New("unknown file type")
	}

	iv := data[9:25]
	encMasterKey := data[25:57]
	encTestString := data[57:61]
	masterKey, err := util.DecryptCFB(encMasterKey, iv, userPW, true)
	if err != nil {
		return nil, err
	}
	testString, err := util.DecryptCFB(encTestString, iv, masterKey, false)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(testString, []byte{116, 114, 117, 101}) {
		return nil, errors.New("failed to decode test string")
	}
	decData, err := util.DecryptCFB(data[61:], iv, masterKey, false)
	if err != nil {
		return nil, err
	}
	head := deserializeHeadCore(decData)
	head.useEncryption = true
	head.masterKey = masterKey
	return head, nil
}
