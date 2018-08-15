package broker

type Publication interface {
	Topic() string
	Message() interface{}
}

type Handler func(Publication)

type Brocker interface {
	Publish(topic string, msg interface{})
	Subscribe(topic string, handler Handler)
}

var (
	Default Brocker
)

func Publish(topic string, msg interface{}) {
	Default.Publish(topic, msg)
}

func Subscribe(topic string, handler Handler) {
	Default.Subscribe(topic, handler)
}
