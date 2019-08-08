package main

import (
	"strconv"
)

// chunk breaks the slice into multiple, smaller slice of a given size
func chunk(slice []string, chunkSize int) [][]string {
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

func parseInt(s string, base int, bitSize int, errMessage string, errCallback func(message string)) (int64, bool) {
	ret, err := strconv.ParseInt(s, base, bitSize)

	if err != nil {
		errCallback(errMessage)
		return 0, false
	}

	return ret, true
}
