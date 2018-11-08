package table

// 路由规则
type RouteRule struct {
	MsgName string
	SvcName string
	Mode    string // auth: 需要授权 pass: 可通过
	MsgID   int
}

// 路由表，包含多条路由规则
type RouteTable struct {
	Rule []*RouteRule
}
