package model

import "sync"

const (
	ConfigPath = "RouteRule"
)

type RouteRule struct {
	MsgName string
	SvcName string
	Mode    string // auth: 需要授权 pass: 可通过
}

type RouteTable struct {
	Rule []*RouteRule
}

var (
	ruleByMsgType      = map[string]*RouteRule{}
	ruleByMsgTypeGuard sync.RWMutex
)

func GetTargetService(msgName string) *RouteRule {

	ruleByMsgTypeGuard.RLock()
	defer ruleByMsgTypeGuard.RUnlock()

	if rule, ok := ruleByMsgType[msgName]; ok {
		return rule
	}

	return nil
}
func ClearRule() {

	ruleByMsgTypeGuard.Lock()
	ruleByMsgType = map[string]*RouteRule{}
	ruleByMsgTypeGuard.Unlock()
}

func AddRouteRule(rule *RouteRule) {
	log.Debugf("Add route rule: %+v", *rule)

	ruleByMsgTypeGuard.Lock()
	ruleByMsgType[rule.MsgName] = rule
	ruleByMsgTypeGuard.Unlock()
}
