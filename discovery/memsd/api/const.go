package memsd

import "errors"

const (
	servicePrefix = "_svcdesc_"
	MaxValueSize  = 512 * 1024
)

var (
	ErrValueNotExists = errors.New("value not exists")
	ErrValueTooLarge  = errors.New("value too large")
)
