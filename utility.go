package main

import (
	"bytes"
	"errors"
)

type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

func (v Version) equals(v2 Version) bool {
	return v.Major == v2.Major && v.Minor == v2.Minor && v.Patch == v2.Patch
}

func newVersion(data []byte) Version {
	v := Version{}
	v.Major = data[0]
	v.Minor = data[1]
	v.Patch = data[2]
	return v
}

type Head struct {
	UseEncryption          bool
	MasterPassword         []byte
	Version                Version
	PageSize               uint16
	PageCount              uint32
	TableListLocation      uint32
	EmptyPagesListLocation uint32
}

func (h *Head) equals(h2 *Head) bool {
	return h.UseEncryption == h2.UseEncryption &&
		bytes.Equal(h.MasterPassword, h2.MasterPassword) &&
		h.Version.equals(h2.Version) &&
		h.PageSize == h2.PageSize &&
		h.PageCount == h2.PageCount &&
		h.TableListLocation == h2.TableListLocation &&
		h.EmptyPagesListLocation == h2.EmptyPagesListLocation
}

func (h *Head) serializeHeadCore() []byte {
	resp := make([]byte, 17)
	resp[0] = h.Version.Major
	resp[1] = h.Version.Minor
	resp[2] = h.Version.Patch
	copy(resp[3:5], uint16toBytes(h.PageSize))
	copy(resp[5:9], uint32toBytes(h.PageCount))
	copy(resp[9:13], uint32toBytes(h.TableListLocation))
	copy(resp[13:17], uint32toBytes(h.EmptyPagesListLocation))
	return resp
}

func deserializeHeadCore(data []byte) *Head {
	h := Head{}
	h.Version = newVersion(data[0:3])
	h.PageSize = bytesToUInt16(data[3:5])
	h.PageCount = bytesToUInt32(data[5:9])
	h.TableListLocation = bytesToUInt32(data[9:13])
	h.EmptyPagesListLocation = bytesToUInt32(data[13:17])
	return &h
}

func (h *Head) SerializeHead(userPW []byte) ([]byte, error) {
	c := h.serializeHeadCore()
	head := make([]byte, 128)
	if !h.UseEncryption {
		copy(head[:5], []byte{103, 111, 68, 66, 00})
		copy(head[5:], c)
		return head, nil
	}
	copy(head[:9], []byte{103, 111, 68, 66, 32, 101, 110, 99, 00})
	iv, err := randomIV(16)
	if err != nil {
		return nil, err
	}
	copy(head[9:25], iv)
	masterKey, err := encryptCFB(h.MasterPassword, iv, userPW)
	if err != nil {
		return nil, err
	}
	copy(head[25:57], masterKey)
	testString, err := encryptCFB(stringtoBytes("true"), iv, h.MasterPassword)
	if err != nil {
		return nil, err
	}
	copy(head[57:61], testString)
	encData, err := encryptCFB(c, iv, h.MasterPassword)
	if err != nil {
		return nil, err
	}
	copy(head[61:], encData)
	return head, nil
}

func DeserializeHead(data []byte, userPW []byte) (*Head, error) {
	if bytes.Equal(data[:5], []byte{103, 111, 68, 66, 00}) {
		return deserializeHeadCore(data[5:]), nil
	}
	if !bytes.Equal(data[:9], []byte{103, 111, 68, 66, 32, 101, 110, 99, 00}) {
		return nil, errors.New("unknown file type")
	}
	iv := data[9:25]
	encKey := data[25:57]
	validation := data[57:61]
	key, err := decryptCFB(encKey, iv, userPW)
	if err != nil {
		return nil, err
	}
	v, err := decryptCFB(validation, iv, key)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(v, stringtoBytes("true")) {
		return nil, errors.New("failed to decode")
	}
	decData, err := decryptCFB(data[61:], iv, key)
	if err != nil {
		return nil, err
	}
	head := deserializeHeadCore(decData)
	head.UseEncryption = true
	head.MasterPassword = key
	return head, nil
}
