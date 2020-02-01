package rbase

import (
	"flag"
)

var (
	FlagAddress      = flag.String("addr", "", "robot login address")
	FlagCount        = flag.Int("count", 1, "robot count")
	FlagCase         = flag.String("case", "", "robot case logic to run")
	FlagShowMsgLog   = flag.Bool("msglog", false, "show msg log")
	FlagFastFastExec = flag.Bool("fastexec", false, "no emulate delay on client")
	FlagRecvTimeOut  = flag.Int("recvtimeout", 5, "recv msg time out in seconds")
)
