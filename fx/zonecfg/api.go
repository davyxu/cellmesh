package zonecfg

import (
	tabtoy "github.com/davyxu/tabtoy/v3/api/golang"
	tabmodel "github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/ulog"
	"sort"
	"strconv"
)

const (
	ZoneConfigFile = "../cfg/ZoneConfig.json"
)

func Load(node, zone string) (err error) {
	currZone = zone
	currNode = node

	ulog.Debugf("Load ZoneConfig zone: %s, file: %s...", zone, ZoneConfigFile)
	err = tabtoy.LoadFromFile(Tab, ZoneConfigFile)

	if err != nil {
		return err
	}

	for _, def := range List() {

		ulog.Infof("ZoneConfig: %s = '%s'", def.Key, def.Value)
	}

	return err
}

func Raw(key string) *ZoneConfig {

	if v, ok := valueByKey.Load(key); ok {
		return v.(*ZoneConfig)
	}

	return nil
}

func List() (ret []*ZoneConfig) {

	valueByKey.Range(func(key, value interface{}) bool {
		ret = append(ret, value.(*ZoneConfig))
		return true
	})

	// 按原导入顺序排序
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Line < ret[j].Line
	})

	return

}

func String(key string) string {

	if v, ok := valueByKey.Load(key); ok {
		return v.(*ZoneConfig).Value
	}

	return ""
}

func Bool(key string) bool {
	v := String(key)
	vv, _ := tabmodel.ParseBool(v)
	return vv
}

func Int(key string) int {
	v := String(key)
	vv, _ := strconv.Atoi(v)
	return vv
}

func Int32(key string) int32 {
	v := String(key)
	vv, _ := strconv.Atoi(v)
	return int32(vv)
}

func Int64(key string) int64 {
	v := String(key)
	vv, _ := strconv.Atoi(v)
	return int64(vv)
}

func Float32(key string) float32 {
	v := String(key)
	vv, _ := strconv.ParseFloat(v, 32)
	return float32(vv)
}
