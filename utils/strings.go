package utils

import (
	"bytes"
)

// StringInSlice checks if the string is present
// in the slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ConcatStrings provides a more performant way
// of concatenating strings than +
func ConcatStrings(strs ...string) string {
	var concatbuffer bytes.Buffer
	for _, str := range strs {
		concatbuffer.WriteString(str)
	}
	return concatbuffer.String()
}
