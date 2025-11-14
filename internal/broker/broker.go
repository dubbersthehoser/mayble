package broker

type Event interface {
	Topic() string
}

type Emiter interface {
	On(topic string, handler func(Event))
	Emit(topic string, any data)
}
