package meshutil

import (
	"flag"
	"github.com/davyxu/cellnet/util"
)

func ApplyFlagFromFile(fs *flag.FlagSet, filename string) error {

	return util.ReadKVFile(filename, func(key, value string) bool {

		// 设置flagm
		fg := fs.Lookup(key)
		if fg != nil {
			log.Infof("ApplyFlagFromFile: %s=%s", key, value)
			fg.Value.Set(value)
		} else {
			log.Errorf("ApplyFlagFromFile: flag not found, %s", key)
		}

		return true
	})
}
