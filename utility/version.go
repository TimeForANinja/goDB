/*
 * Version structure and helper functionss
 */

package utility

// Version represents a version in the XX.XX.XX format
type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

// Equals compares whether two Versions are identical
func (v Version) Equals(v2 Version) bool {
	return v.Major == v2.Major && v.Minor == v2.Minor && v.Patch == v2.Patch
}

// ToBytes returns the Version as byte array
func (v Version) ToBytes() []byte {
	return []byte{v.Major, v.Minor, v.Patch}
}

// NewVersionFromBytes takes at (3) byte array to create a new Version
func NewVersionFromBytes(data []byte) Version {
	v := Version{}
	v.Major = data[0]
	v.Minor = data[1]
	v.Patch = data[2]
	return v
}

// NewVersion takes 3 uint8's to create a new Version
func NewVersion(major, minor, patch uint8) Version {
	return Version{major, minor, patch}
}
