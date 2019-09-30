package memsd

import "errors"

const (
	MaxValueSize = 512 * 1024
)

var (
	ErrValueNotExists = errors.New("value not exists")
	ErrValueTooLarge  = errors.New("value too large")
)
