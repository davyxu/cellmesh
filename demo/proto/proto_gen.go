// Auto generated by github.com/davyxu/cellmesh/protogen
// DO NOT EDIT!

package proto

import (
	"fmt"
	"reflect"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

type VerifyREQ struct {
	Token string
}

type VerifyACK struct {
	Status int32
}

func (self *VerifyREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *VerifyACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*VerifyREQ)(nil)).Elem(),
		ID:    23773,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*VerifyACK)(nil)).Elem(),
		ID:    51140,
	})
}
