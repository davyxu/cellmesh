package service

import "github.com/davyxu/cellmesh/discovery"

type QueryServiceOp int

const (
	QueryServiceOp_NextFilter QueryServiceOp = iota
	QueryServiceOp_NextDesc
	QueryServiceOp_End
)

type QueryResult interface{}

// 返回值含义:
// 1. true等效于QueryServiceOp_NextFilter,转到下一个内层循环
// 2. false等效于QueryServiceOp_NextDesc, 转到下一个外层循环
// 3. QueryServiceOp_End: 终止所有遍历循环
// 4. Filter中将类型转为QueryResult,则在QueryService函数返回
type FilterFunc func(*discovery.ServiceDesc) interface{}

// 根据给定的查询服务名,将结果经过各种过滤器处理后输出
func QueryService(svcName string, filterList ...FilterFunc) (ret interface{}) {

	for _, desc := range discovery.Default.Query(svcName) {

		for _, filter := range filterList {

			if filter == nil {
				continue
			}

			op := filter(desc)

			switch raw := op.(type) {
			case QueryServiceOp:
				switch raw {
				case QueryServiceOp_NextFilter:
				case QueryServiceOp_NextDesc:
					goto NextDesc
				case QueryServiceOp_End:
					return
				}
			case bool:
				if !raw {
					goto NextDesc
				}
			case QueryResult:
				ret = raw
			default:
				panic("unknown filter result")
			}
		}

	NextDesc:
	}

	return
}

// 匹配指定的服务组,服务组空时,匹配所有
func Filter_MatchSvcGroup(svcGroup string) FilterFunc {

	return func(desc *discovery.ServiceDesc) interface{} {

		if svcGroup == "" {
			return true
		}

		return desc.GetMeta("SvcGroup") == svcGroup
	}
}

// 匹配指定的服务ID
func Filter_MatchSvcID(svcid string) FilterFunc {

	return func(desc *discovery.ServiceDesc) interface{} {

		if desc.ID == svcid {
			return QueryResult(desc)
		}

		return true
	}
}

// 匹配指定的规则,一般由命令行指定
func Filter_MatchRule(rules []MatchRule) FilterFunc {

	return func(desc *discovery.ServiceDesc) interface{} {

		// 任意规则满足即可
		for _, rule := range rules {

			if matchTarget(&rule, desc) {
				return true
			}
		}

		return false
	}

}
