package meshutil

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestUUID(t *testing.T) {

	for i := 0; i < 10; i++ {
		u := uuid.NewV4()
		t.Log(u.String(), len(u.String()))
	}

}
