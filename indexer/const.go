package indexer

import "errors"

const (
	SEEK_BEGINNING  = 0
	SIZE_OF_INT64   = 8
	SIZE_OF_NEWLINE = 1
)

var (
	ouchErr   = errors.New("Negative line")
	offsetErr = errors.New("Offset error")
)
