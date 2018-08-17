package model

import "sync"

const (
	ConfigPath = "config/agent/route_rule"
)

// 路由规则
type RouteRule struct {
	MsgName string
	SvcName string
	Mode    string // auth: 需要授权 pass: 可通过
}

// 路由表，包含多条路由规则
type RouteTable struct {
	Rule []*RouteRule
}

var (
	// 消息名映射路由规则
	ruleByMsgName      = map[string]*RouteRule{}
	ruleByMsgNameGuard sync.RWMutex
)

// 消息名取路由规则
func GetTargetService(msgName string) *RouteRule {

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
	ruleByMsgName = map[string]*RouteRule{}
	ruleByMsgNameGuard.Unlock()
}

// 添加路由规则
func AddRouteRule(rule *RouteRule) {
	log.Debugf("Add route rule: %+v", *rule)

	ruleByMsgNameGuard.Lock()
	ruleByMsgName[rule.MsgName] = rule
	ruleByMsgNameGuard.Unlock()
}
