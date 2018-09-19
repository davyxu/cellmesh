package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"strings"
)

type matchRule struct {
	SvcName    string
	TargetNode string
}

var (
	matchRules []matchRule
)

func matchTarget(node string, desc *discovery.ServiceDesc) bool {

	// Tags中保存服务所在的节点
	for _, sdTag := range desc.Tags {
		if sdTag == node {
			return true
		}
	}

	return false
}

func MatchService(svcName string, desclist []*discovery.ServiceDesc) (ret []*discovery.ServiceDesc) {
	return rawMatch(matchRules, svcName, desclist)
}

// 获取要匹配节点名(连接用)
func rawMatch(rules []matchRule, svcName string, desclist []*discovery.ServiceDesc) (ret []*discovery.ServiceDesc) {

	if len(desclist) == 0 {
		return
	}

	// 优先精确匹配

	var anyMatch bool
	for _, rule := range rules {
		if rule.SvcName == svcName {

			anyMatch = true

			for _, sd := range desclist {

				if matchTarget(rule.TargetNode, sd) {
					ret = append(ret, sd)
				}
			}
		}
	}

	// 精确匹配到，跳过默认方式
	if anyMatch {
		return
	}

	// 默认方式
	for _, rule := range rules {
		if rule.SvcName == "" {

			anyMatch = true

			for _, sd := range desclist {

				if matchTarget(rule.TargetNode, sd) {
					ret = append(ret, sd)
				}
			}
		}
	}

	return
}

func parseMatchRule(rule string) (ret []matchRule) {

	for _, ruleStr := range strings.Split(rule, "|") {

		ruleStr = strings.TrimSpace(ruleStr)
		rulePairs := strings.Split(ruleStr, ":")

		var rule matchRule
		if len(rulePairs) == 2 {
			rule.SvcName = rulePairs[0]
			rule.TargetNode = rulePairs[1]
		} else {
			rule.TargetNode = ruleStr
		}

		ret = append(ret, rule)
	}

	return
}
