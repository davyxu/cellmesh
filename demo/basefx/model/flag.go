package fxmodel

import (
	"github.com/davyxu/cellmesh/service"
)

var (
	FlagSelfGroup       = service.CommandLine.Bool("forceselfgroup", false, "Force match curr svcgroup")
	FlagCommunicateType = service.CommandLine.String("commtype", "tcp", "Communicate type, tcp or ws")
)
