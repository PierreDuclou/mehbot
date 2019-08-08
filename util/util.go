package util

import (
	"strconv"
)

// Chunk breaks the slice into multiple, smaller slice of a given size
func Chunk(slice []string, chunkSize int) [][]string {
	var chunks [][]string

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// ParseInt is a helper wrapping the strconv.ParseInt function and allowing to throw error messages using a callback
func ParseInt(s string, base int, bitSize int, errMessage string, errCallback func(message string)) (int64, bool) {
	ret, err := strconv.ParseInt(s, base, bitSize)

	if err != nil {
		errCallback(errMessage)
		return 0, false
	}

	return ret, true
}
