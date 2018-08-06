package router

import (
	"github.com/davyxu/cellnet"
	"reflect"
)

func QuerySerivceByMsgType(msgType reflect.Type) (string, bool) {

	meta := cellnet.MessageMetaByType(msgType)
	if raw, ok := meta.GetContext("service"); ok {
		return raw.(string), true
	}

	return "", false
}
