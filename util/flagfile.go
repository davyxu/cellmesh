package meshutil

import (
	"flag"
	"github.com/davyxu/cellnet/util"
	"github.com/davyxu/ulog"
)

func ApplyFlagFromFile(fs *flag.FlagSet, filename string) error {

	return util.ReadKVFile(filename, func(key, value string) bool {

		// 设置flagm
		fg := fs.Lookup(key)
		if fg != nil {
			ulog.Infof("ApplyFlagFromFile: %s=%s", key, value)
			fg.Value.Set(value)
		} else {
			ulog.Errorf("ApplyFlagFromFile: flag not found, %s", key)
		}

		return true
	})
}
