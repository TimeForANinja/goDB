package goDB

import (
	"testing"
)

func TestPage(t *testing.T) {
	myPageSize := 128
	myData := make([]byte, 64)
	copy(myData, []byte{255, 254, 253, 252, 251, 250, 249, 248, 247, 246, 245, 244, 243, 242, 241, 240, 0, 0, 0, 0})
	myPageSpezific := make([]byte, 35)
	copy(myPageSpezific, []byte{69})
	myPageHead := pageHead{
		pageType:     1,
		nextPage:     1337,
		prevPage:     420,
		firstItem:    10,
		endTrim:      69,
		pageSpezific: myPageSpezific,
	}
	myPage := page{
		index:    2,
		pageHead: &myPageHead,
		data:     myData,
	}

	// serialize page
	testPageSerial, err := myPage.serializePage(uint16(myPageSize), false, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(testPageSerial) != myPageSize {
		t.Error("testPageSerial invalid length", testPageSerial, myPageSize)
	}
	// deserialize again
	decodedHead, decodedBytes, err := deserializePage(testPageSerial, false, nil)
	if err != nil {
		t.Error(err)
		return
	}
	testPageDecoded := page{
		index:    2,
		pageHead: decodedHead,
		data:     decodedBytes,
	}
	// compare
	if !myPage.equals(&testPageDecoded) {
		t.Error("objects not equal", myPage, testPageDecoded)
	}
}
