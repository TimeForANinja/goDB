package main

import (
	"bytes"
	"errors"

	util "github.com/timeforaninja/goDB/utility"
)

type head struct {
	UseEncryption          bool
	MasterPassword         []byte
	Version                util.Version
	PageSize               uint16
	PageCount              uint32
	TableListLocation      uint32
	EmptyPagesListLocation uint32
}

func (h *head) equals(h2 *head) bool {
	return h.UseEncryption == h2.UseEncryption &&
		bytes.Equal(h.MasterPassword, h2.MasterPassword) &&
		h.Version.Equals(h2.Version) &&
		h.PageSize == h2.PageSize &&
		h.PageCount == h2.PageCount &&
		h.TableListLocation == h2.TableListLocation &&
		h.EmptyPagesListLocation == h2.EmptyPagesListLocation
}

func (h *head) serializeHeadCore() []byte {
	resp := make([]byte, 17)
	copy(resp[:3], h.Version.ToBytes())
	copy(resp[3:5], util.Uint16toBytes(h.PageSize))
	copy(resp[5:9], util.Uint32toBytes(h.PageCount))
	copy(resp[9:13], util.Uint32toBytes(h.TableListLocation))
	copy(resp[13:17], util.Uint32toBytes(h.EmptyPagesListLocation))
	return resp
}

func deserializeHeadCore(data []byte) *head {
	h := head{}
	h.Version = util.NewVersionFromBytes(data[0:3])
	h.PageSize = util.BytesToUInt16(data[3:5])
	h.PageCount = util.BytesToUInt32(data[5:9])
	h.TableListLocation = util.BytesToUInt32(data[9:13])
	h.EmptyPagesListLocation = util.BytesToUInt32(data[13:17])
	return &h
}

func (h *head) serializeHead(userPW []byte) ([]byte, error) {
	c := h.serializeHeadCore()
	head := make([]byte, 128)
	if !h.UseEncryption {
		copy(head[:5], []byte{103, 111, 68, 66, 00})
		copy(head[5:], c)
		return head, nil
	}
	copy(head[:9], []byte{103, 111, 68, 66, 32, 101, 110, 99, 00})
	masterKey, iv, err := util.NewEncryptCFB(h.MasterPassword, userPW)
	if err != nil {
		return nil, err
	}
	copy(head[9:25], iv)
	copy(head[25:57], masterKey)
	testString, err := util.EncryptCFB(util.StringtoBytes("true"), iv, h.MasterPassword)
	if err != nil {
		return nil, err
	}
	copy(head[57:61], testString)
	encData, err := util.EncryptCFB(c, iv, h.MasterPassword)
	if err != nil {
		return nil, err
	}
	copy(head[61:], encData)
	return head, nil
}

func deserializeHead(data []byte, userPW []byte) (*head, error) {
	if bytes.Equal(data[:5], []byte{103, 111, 68, 66, 00}) {
		return deserializeHeadCore(data[5:]), nil
	}
	if !bytes.Equal(data[:9], []byte{103, 111, 68, 66, 32, 101, 110, 99, 00}) {
		return nil, errors.New("unknown file type")
	}
	iv := data[9:25]
	encKey := data[25:57]
	validation := data[57:61]
	key, err := util.DecryptCFB(encKey, iv, userPW)
	if err != nil {
		return nil, err
	}
	v, err := util.DecryptCFB(validation, iv, key)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(v, util.StringtoBytes("true")) {
		return nil, errors.New("failed to decode")
	}
	decData, err := util.DecryptCFB(data[61:], iv, key)
	if err != nil {
		return nil, err
	}
	head := deserializeHeadCore(decData)
	head.UseEncryption = true
	head.MasterPassword = key
	return head, nil
}
