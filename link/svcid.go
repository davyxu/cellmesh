package link

import (
	"fmt"
	"github.com/davyxu/cellmesh/fx"
)

// 全局唯一的svcid 格式:  svcName@ip
// 一台机器开多套时, 需要配合FlagFile指定每组的GroupName, 或者手动参数指定

// 构造指定服务的ID
func MakeSvcID(svcName string) string {
	return fmt.Sprintf("%s#%d@%s", svcName, fx.SvcIndex, fx.SvcGroup)
}

// 获得本进程的服务id
func GetLocalSvcID() string {
	return MakeSvcID(fx.ProcName)
}
