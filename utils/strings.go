package utils

import (
	"bytes"
	"io"
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

// ReadAllString returns the contents of the
// io.Reader stringified
func ReadAllString(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String()
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
