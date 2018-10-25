package model

import (
	"github.com/davyxu/cellmesh/demo/table"
	"sync"
)

const (
	ConfigPath = "cm_demo/config/agent/route_rule"
)

var (
	// 消息名映射路由规则
	ruleByMsgName      = map[string]*table.RouteRule{}
	ruleByMsgNameGuard sync.RWMutex
)

// 消息名取路由规则
func GetTargetService(msgName string) *table.RouteRule {

	ruleByMsgNameGuard.RLock()
	defer ruleByMsgNameGuard.RUnlock()

	if rule, ok := ruleByMsgName[msgName]; ok {
		return rule
	}

	return nil
}

// 清除所有规则
func ClearRule() {

	ruleByMsgNameGuard.Lock()
	ruleByMsgName = map[string]*table.RouteRule{}
	ruleByMsgNameGuard.Unlock()
}

// 添加路由规则
func AddRouteRule(rule *table.RouteRule) {

	ruleByMsgNameGuard.Lock()
	ruleByMsgName[rule.MsgName] = rule
	ruleByMsgNameGuard.Unlock()
}
