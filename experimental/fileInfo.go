package main

/*
 * utility for the meta information in the first 128 bytes
 */
type fileInfo struct {
	version              version
	databasePageSize     int16
	fileChangeCounter    int32
	databasePages        int32
	defaultPageCacheSize int32
}

func (f fileInfo) toBytes() []byte {
	resp := make([]byte, 17)
	resp[0] = byte(f.version.major)
	resp[1] = byte(f.version.minor)
	resp[2] = byte(f.version.patch)
	copy(resp[3:5], int16toBytes(f.databasePageSize))
	copy(resp[5:9], int32toBytes(f.fileChangeCounter))
	copy(resp[9:13], int32toBytes(f.databasePages))
	copy(resp[13:17], int32toBytes(f.defaultPageCacheSize))
	return resp
}

func parseFileInfo(data []byte) fileInfo {
	f := fileInfo{}
	f.version = version{
		major: int8(data[0]),
		minor: int8(data[1]),
		patch: int8(data[2]),
	}
	f.databasePageSize = bytesToInt16(data[3:5])
	f.fileChangeCounter = bytesToInt32(data[5:9])
	f.databasePages = bytesToInt32(data[9:13])
	f.defaultPageCacheSize = bytesToInt32(data[13:17])
	return f
}
func defaultFileInfo() fileInfo {
	return fileInfo{
		version:              version{major: 0, minor: 0, patch: 1},
		databasePageSize:     4096,
		fileChangeCounter:    0,
		databasePages:        0,
		defaultPageCacheSize: 0,
	}
}
