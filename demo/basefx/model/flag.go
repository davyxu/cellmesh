package fxmodel

import "flag"

var (
	FlagSelfGroup       = flag.Bool("forceselfgroup", false, "Force match curr svcgroup")
	FlagCommunicateType = flag.String("commtype", "tcp", "Communicate type, tcp or ws")
)
