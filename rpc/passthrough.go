package rpc

import (
	"fmt"
	"github.com/davyxu/cellmesh/proto"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"reflect"
	"strings"
)

func encodePassthrough(passthrough interface{}) ([]byte, string, error) {

	var typeName string
	switch pvalue := passthrough.(type) {
	case int64:
		passthrough = &proto.PassThroughWrap{
			Int64: pvalue,
		}
		typeName = "pt.int64"
	case int32:
		passthrough = &proto.PassThroughWrap{
			Int32: pvalue,
		}
		typeName = "pt.int32"
	case string:
		passthrough = &proto.PassThroughWrap{
			Str: pvalue,
		}
		typeName = "pt.string"
	case float32:
		passthrough = &proto.PassThroughWrap{
			Float32: pvalue,
		}
		typeName = "pt.float32"
	}

	data, meta, err := codec.EncodeMessage(passthrough, nil)

	if err != nil {
		return nil, "", err
	}

	if typeName != "" {
		return data, typeName, nil
	}

	return data, meta.FullName(), nil
}

var ptwMeta *cellnet.MessageMeta

func init() {
	ptwMeta = cellnet.MessageMetaByType(reflect.TypeOf((*proto.PassThroughWrap)(nil)).Elem())
}

func decodePassthrough(data []byte, dataType string) (pt interface{}, err error) {

	if strings.HasPrefix(dataType, "pt.") {

		msg := &proto.PassThroughWrap{}
		err = ptwMeta.Codec.Decode(data, msg)

		if err != nil {
			return nil, err
		}

		switch dataType {
		case "pt.int64":
			pt = msg.Int64
		case "pt.int32":
			pt = msg.Int32
		case "pt.string":
			pt = msg.Str
		case "pt.float32":
			pt = msg.Float32
		}

		return

	} else {

		// 获取消息元信息
		meta := cellnet.MessageMetaByFullName(dataType)

		// 消息没有注册
		if meta == nil {
			return nil, fmt.Errorf("type not found: %s", dataType)
		}

		// 创建消息
		pt = meta.NewType()

		// 从字节数组转换为消息
		err = meta.Codec.Decode(data, pt)

		if err != nil {
			return nil, err
		}

		return
	}
}
