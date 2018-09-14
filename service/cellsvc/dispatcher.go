package cellsvc

import "github.com/davyxu/cellmesh/service"

type evDispatcher struct {
	dispatcherFunc service.EventFunc
	procName       string // "tcp.ltv"
}

func (self *evDispatcher) SetProcessor(name string) {
	self.procName = name
}

func (self *evDispatcher) SetEventCallback(dis service.EventFunc) {
	self.dispatcherFunc = dis
}

func (self *evDispatcher) Invoke(event service.Event) {

	if self.dispatcherFunc != nil {
		self.dispatcherFunc(event)
	}
}
