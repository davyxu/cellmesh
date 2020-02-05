package model

import (
	"fmt"
	"github.com/davyxu/cellmesh/svc/robot/rbase"
	"github.com/davyxu/ulog"
)

type Robot struct {
	rbase.Messenger // 网络

	// 基础结构
	ID string
	// 账号

	AgentSvcID   string
	AgentAddress string // 网关地址
	LoginToken   string // login返回的token
	GameToken    string // 断线重连的token

	state string
}

func (self *Robot) SetState(state string) {
	self.state = state
	AddState(state)
	ulog.Infof("%s SetState: %s", self.AccountName(), state)

	self.Sleep()
}

func (self *Robot) State() string {
	return self.state
}

func (self *Robot) AccountName() string {
	return fmt.Sprintf("r%s", self.ID)
}

func (self *Robot) RunFlow(flow func(*Robot)) {

	go flow(self)
}

func NewRobot(id string) *Robot {

	self := &Robot{
		ID: id,
	}

	self.Init()

	return self
}
