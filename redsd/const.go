package redsd

import (
	"encoding/json"
	"github.com/davyxu/ulog"
)

const (
	key_NodeReg  = "sd_reg:"  // hash:svcid value:1 按名称分类节点的注册节点
	key_Node     = "sd_node:" // 实际详细信息, 服务自己维护心跳
	key_NodeLock = "sd_lock:" // node操作锁, 关联nodereg
	key_NodeDesc = "Desc"     // 节点服务发现信息
	key_NodeVer  = "Ver"      // 节点信息版本号
	key_KV       = "sd_kv:"   // KV存储
)

const (
	nodeKeyTimeoutSec       = 7 // node信息超时时间
	nodeKeyCheckIntervalSec = 3 //  多久检查一次node超时
)

func marshalDesc(desc *NodeDesc) []byte {
	data, err := json.Marshal(desc)
	if err != nil {
		ulog.Errorf("marshal redisDiscovery desc failed, %s", err)
		return nil
	}

	return data
}

func unmarshalDesc(data []byte) (ret *NodeDesc) {
	ret = NewDesc()
	err := json.Unmarshal(data, ret)
	if err != nil {
		ulog.Errorf("unmarshal redisDiscovery desc failed, %s", err)
	}

	return
}
