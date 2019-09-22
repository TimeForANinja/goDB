package main

/*
 * everything about a page in the database
 */
type tablePage struct {
	indexOfTable            int64
	numberOfColumns         int64
	indexOfColumnDefPage    int64
	indexOfFirstContentPage int64
	pageCount               int64
	rowCount                int64
	tableNameLength         int64
	tableName               string
}

func readPage(db *Database) *tablePage {
	var data []byte
	if db.encryption {
		data = readEncryptedPage(db)
	} else {
		data = readUnencryptedPage(db)
	}
	return parsePage(data)
}
func readEncryptedPage(db *Database) []byte {
	// read db.fileinfo.databasePageSize bytes
	// use last aes.BlockSize bytes as iv
	// return the decrypted data
	return nil
}
func readUnencryptedPage(db *Database) []byte {
	// read db.fileinfo.databasePageSize bytes
	return nil
}
func parsePage([]byte) *tablePage {
	return nil
}
