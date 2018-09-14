package fxmodel

var (
	Node       string   // svcid的尾缀，与名字组合后全局唯一
	MatchNodes []string // 去匹配其他节点
)

func GetSvcID(svcName string) string {
	return svcName + "@" + Node
}
