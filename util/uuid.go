package meshutil

import (
	"github.com/satori/go.uuid"
)

func GenID() string {
	id := uuid.NewV4()
	return id.String()
}
