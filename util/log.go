package meshutil

import (
	"github.com/davyxu/golog"
	"strconv"
	"strings"
)

var log = golog.New("meshutil")

func sizeLevel(sizeStr, levelStr string, multi int, sizePtr *int, errPtr *error) bool {
	if strings.HasSuffix(sizeStr, levelStr) {
		size, err := strconv.Atoi(strings.TrimSuffix(sizeStr, levelStr))

		if err != nil {
			*errPtr = err
			return true
		}

		*sizePtr = size * multi
		return true
	}

	return false
}

func ParseSizeString(sizeStr string) (size int, err error) {

	sizeStr = strings.TrimSpace(sizeStr)
	sizeStr = strings.ToUpper(sizeStr)

	if sizeLevel(sizeStr, "M", 1024*1024, &size, &err) {
		return
	}

	if sizeLevel(sizeStr, "K", 1024, &size, &err) {
		return
	}

	if sizeLevel(sizeStr, "G", 1024*1024*1024, &size, &err) {
		return
	}

	return strconv.Atoi(sizeStr)
}
