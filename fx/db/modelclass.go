package db

import (
	"reflect"
	"strconv"
	"strings"
)

type ModelClass interface {
	DBKey() string // redis外层key的前缀
}

func genKey(key string, rawKeyID interface{}) string {
	var sb strings.Builder
	sb.WriteString(key)
	sb.WriteString(":")

	switch keyID := rawKeyID.(type) {
	case int64:
		sb.WriteString(strconv.Itoa(int(keyID)))
	case string:
		sb.WriteString(keyID)
	default:
		panic("unsupport keyid type: " + reflect.TypeOf(keyID).Name())
	}

	return sb.String()
}
