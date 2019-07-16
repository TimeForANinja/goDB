package main

/*
 Database Flags
*/
type Flags struct {
	UNIQUE         bool
	NOT_NULL       bool
	AUTO_INCREMENT bool
	PRIMARY_KEY    bool
}

func FlagsFromStrings(s []string) Flags {
	output := Flags{false, false, false, false}
	if contains(s, "UNIQUE") {
		output.UNIQUE = true
	}
	if contains(s, "NOT_NULL") {
		output.NOT_NULL = true
	}
	if contains(s, "AUTO_INCREMENT") {
		output.AUTO_INCREMENT = true
	}
	if contains(s, "PRIMARY_KEY") {
		output.PRIMARY_KEY = true
	}
	return output
}
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
func FlagsFromByte(b byte) Flags {
	return Flags{
		UNIQUE:         uint8(b)&1 == 1,
		NOT_NULL:       uint8(b)&2 == 1,
		AUTO_INCREMENT: uint8(b)&4 == 1,
		PRIMARY_KEY:    uint8(b)&8 == 1,
	}
}
func (f Flags) toByte() byte {
	var output uint8
	if f.UNIQUE {
		output += 1
	}
	if f.NOT_NULL {
		output += 2
	}
	if f.AUTO_INCREMENT {
		output += 4
	}
	if f.PRIMARY_KEY {
		output += 8
	}
	return byte(output)
}
