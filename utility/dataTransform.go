/*
 * type conversion functions
 */

package utility

import "encoding/binary"

// Uint16toBytes converts a uint16 to a byte array
func Uint16toBytes(num uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, num)
	return b
}

// Uint32toBytes converts a uint32 to a byte array
func Uint32toBytes(num uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, num)
	return b
}

// StringtoBytes converts a string to a (utf8) byte array
func StringtoBytes(s string) []byte {
	return []byte(s)
}

// BytesToUInt16 converts a byte array to uint16
func BytesToUInt16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

// BytesToUInt32 converts a byte array to uint32
func BytesToUInt32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

// BytesToString converts a byte array to a (utf8) string
func BytesToString(bytes []byte) string {
	return string(bytes)
}

/* int conversion
func Int16toBytes(num int16) []byte {
	return []byte{
		byte(uint >> 8),
		byte(uint & 255),
	}
}

func Int32toBytes(num int32) []byte {
	return []byte{
		byte((int >> 24) & 255),
		byte((int >> 16) & 255),
		byte((int >> 8) & 255),
		byte(int & 255),
	}
}

func BytesToInt16(bytes []byte) int16 {
	return (uint16(bytes[0]) << 8) + uint16(bytes[1])
}

func BytesToInt32(bytes []byte) int32 {
	return (uint32(bytes[2]) << 24) + (uint32(bytes[2]) << 16) + (uint32(bytes[2]) << 8) + uint32(bytes[3])
}
*/
