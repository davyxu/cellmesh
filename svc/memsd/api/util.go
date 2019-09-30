package memsd

import (
	"fmt"
	"github.com/davyxu/cellmesh/svc/memsd/proto"
)

func codeToError(code sdproto.ResultCode) error {

	switch code {
	case sdproto.ResultCode_Result_OK:
		return nil
	case sdproto.ResultCode_Result_NotExists:
		return ErrValueNotExists
	}

	return fmt.Errorf("error %s", code.String())
}
