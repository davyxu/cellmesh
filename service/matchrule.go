package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"regexp"
	"strings"
)

type MatchRule struct {
	SvcName string
	Target  string

	nodeExp *regexp.Regexp
}

func (self *MatchRule) MatchNode(node string) bool {

	if self.nodeExp == nil {
		exp, err := regexp.Compile(self.Target)
		if err != nil {
			return false
		}

		self.nodeExp = exp
	}

	return self.nodeExp.MatchString(node)
}

func matchTarget(rule *MatchRule, desc *discovery.ServiceDesc) bool {

	return rule.MatchNode(desc.GetMeta("SvcGroup"))
}

// 获取要匹配节点名(连接用)
func MatchService(rules []MatchRule, svcName string, desclist []*discovery.ServiceDesc) (ret []*discovery.ServiceDesc) {

	if len(desclist) == 0 {
		return
	}

	// 优先精确匹配

	var anyMatch bool
	for _, rule := range rules {
		if rule.SvcName == svcName {

			anyMatch = true

			for _, sd := range desclist {

				if matchTarget(&rule, sd) {
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

			for _, sd := range desclist {

				if matchTarget(&rule, sd) {
					ret = append(ret, sd)
				}
			}
		}
	}

	return
}

func ParseMatchRule(rule string) (ret []MatchRule) {

	for _, ruleStr := range strings.Split(rule, "|") {

		ruleStr = strings.TrimSpace(ruleStr)
		rulePairs := strings.Split(ruleStr, ":")

		var rule MatchRule
		if len(rulePairs) == 2 {
			rule.SvcName = rulePairs[0]
			rule.Target = rulePairs[1]
		} else {
			if ruleStr == "" {
				continue
			}

			rule.Target = ruleStr
		}

		ret = append(ret, rule)
	}

	return
}
