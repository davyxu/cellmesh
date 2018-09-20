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

	rules := ParseMatchRule("game:a|b")

	t.Log(MatchService(rules, "game", makeSD("b", "a")))

	t.Log(MatchService(rules, "hub", makeSD("b", "c")))
}
