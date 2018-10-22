package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"strings"
)

type MatchRule struct {
	Target string
}

func matchTarget(rule *MatchRule, desc *discovery.ServiceDesc) bool {

	return rule.Target == desc.GetMeta("SvcGroup")
}
func ParseMatchRule(rule string) (ret []MatchRule) {

	for _, ruleStr := range strings.Split(rule, "|") {
		var rule MatchRule
		rule.Target = ruleStr
		ret = append(ret, rule)
	}

	return
}
