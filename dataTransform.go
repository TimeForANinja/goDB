package main

import "encoding/binary"

/*
 * data transformation utilities
 */
func uint16toBytes(num uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, num)
	return b
	//	return []byte{
	//		byte(uint >> 8),
	//		byte(uint & 255),
	//	}
}
func uint32toBytes(num uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, num)
	return b
	//	return []byte{
	//		byte((int >> 24) & 255),
	//		byte((int >> 16) & 255),
	//		byte((int >> 8) & 255),
	//		byte(int & 255),
	//	}
}
func stringtoBytes(s string) []byte {
	return []byte(s)
}
func bytesToUInt16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
	// return (uint16(bytes[0]) << 8) + uint16(bytes[1])
}
func bytesToUInt32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
	// return (uint32(bytes[2]) << 24) + (uint32(bytes[2]) << 16) + (uint32(bytes[2]) << 8) + uint32(bytes[3])
}
func bytesToString(bytes []byte) string {
	return string(bytes)
}
