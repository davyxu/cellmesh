package service

import (
	"github.com/davyxu/cellmesh/discovery"
	"testing"
)

func makeSD(svcNodes ...string) (ret []*discovery.ServiceDesc) {

	for _, node := range svcNodes {
		ret = append(ret, &discovery.ServiceDesc{
			Tags: []string{node},
		})
	}

	return
}

func TestMatchRule(t *testing.T) {

	rules := parseMatchRule("game:a|b")

	t.Log(rawMatch(rules, "game", makeSD("b", "a")))

	t.Log(rawMatch(rules, "hub", makeSD("b", "c")))
}
