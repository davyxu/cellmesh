package broker

import "sync"

type localBroker struct {
	handlerByTopic      map[string][]Handler
	handlerByTopicGuard sync.RWMutex
}

type pubCtx struct {
	topic string
	msg   interface{}
}

func (self *pubCtx) Topic() string {
	return self.topic
}

func (self *pubCtx) Message() interface{} {
	return self.msg
}

func (self *localBroker) Publish(topic string, msg interface{}) {

	var list []Handler
	self.handlerByTopicGuard.RLock()

	list, _ = self.handlerByTopic[topic]

	self.handlerByTopicGuard.RUnlock()

	ctx := &pubCtx{topic, msg}

	for _, h := range list {
		h(ctx)
	}
}

func (self *localBroker) Subscribe(topic string, handler Handler) {

	self.handlerByTopicGuard.Lock()

	list, _ := self.handlerByTopic[topic]
	list = append(list, handler)
	self.handlerByTopic[topic] = list

	self.handlerByTopicGuard.Unlock()
}

func NewLocalBroker() Brocker {
	return &localBroker{
		handlerByTopic: make(map[string][]Handler),
	}
}
