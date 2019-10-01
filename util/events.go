package meshutil

type EventHandler func(args ...interface{})

type EventHandlerSet []EventHandler

func (self *EventHandlerSet) Add(h EventHandler) {
	*self = append(*self, h)
}

func (self *EventHandlerSet) Invoke(args ...interface{}) {
	for _, h := range *self {
		h(args...)
	}
}
