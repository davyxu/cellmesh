package memsd

import (
	"fmt"
	"github.com/davyxu/cellmesh/discovery/memsd/proto"
)

func codeToError(code proto.ResultCode) error {

	switch code {
	case proto.ResultCode_Result_OK:
		return nil
	case proto.ResultCode_Result_NotExists:
		return ErrValueNotExists
	}

	return fmt.Errorf("unknown error %d", code)
}
