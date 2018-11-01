package service

import (
	"flag"
	"github.com/davyxu/cellnet/util"
	"strings"
)

func ApplyFlagFromFile(filename string) error {

	return util.ReadFileLines(filename, func(line string) bool {

		line = strings.TrimSpace(line)

		// 注释
		if strings.HasPrefix(line, "#") {
			return true
		}

		// 等号切分KV
		pairs := strings.Split(line, "=")
		if len(pairs) == 2 {

			key := pairs[0]
			value := pairs[1]

			// 设置flag
			fg := flag.Lookup(key)
			if fg != nil {
				log.Infof("ApplyFlagFromFile: %s=%s", key, value)
				fg.Value.Set(value)
			} else {
				log.Errorf("ApplyFlagFromFile: flag not found, %s", key)
			}

		}

		return true
	})
}
