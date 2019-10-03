package meshutil

type EventFunc func(args ...interface{})

type EventFuncSet []EventFunc

func (self *EventFuncSet) Add(h EventFunc) {
	*self = append(*self, h)
}

func (self *EventFuncSet) Invoke(args ...interface{}) {
	for _, h := range *self {
		h(args...)
	}
}
