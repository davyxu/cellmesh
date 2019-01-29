package fxmodel

import (
	"github.com/davyxu/cellmesh/service"
)

var (
	FlagCommunicateType = service.CommandLine.String("commtype", "tcp", "Communicate type, tcp or ws")
)
