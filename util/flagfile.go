package meshutil

import (
	"flag"
	"github.com/davyxu/cellnet/util"
)

func ApplyFlagFromFile(filename string) error {

	return util.ReadKVFile(filename, func(key, value string) bool {

		// 设置flag
		fg := flag.Lookup(key)
		if fg != nil {
			log.Infof("ApplyFlagFromFile: %s=%s", key, value)
			fg.Value.Set(value)
		} else {
			log.Errorf("ApplyFlagFromFile: flag not found, %s", key)
		}

		return true
	})
}
