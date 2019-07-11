package main

/*
 * data transformation utilities
 */
func int16toBytes(int int16) []byte {
	return []byte{
		byte(int >> 8),
		byte(int & 255),
	}
}
func int32toBytes(int int32) []byte {
	return []byte{
		byte((int >> 24) & 255),
		byte((int >> 16) & 255),
		byte((int >> 8) & 255),
		byte(int & 255),
	}
}
func bytesToInt16(bytes []byte) int16 {
	return (int16(bytes[0]) << 8) + int16(bytes[1])
}
func bytesToInt32(bytes []byte) int32 {
	return (int32(bytes[2]) << 24) + (int32(bytes[2]) << 16) + (int32(bytes[2]) << 8) + int32(bytes[3])
}
