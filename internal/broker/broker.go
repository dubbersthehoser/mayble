package broker




type Broker struct {
	event map[string][]func()
	request map[string]func(func(any))
}

func NewBroker() *Broker {
	return &Broker{
		direct: make(map[string[]Handler),
	}
}

func (b *Broker) OnRequest(key string, f func(func(any))) {
	b.request[key] = f
}

func (b *Broker) Request(key string, f func(any)) error {
	cmd, ok := b.request[key]
	if !ok {
		return errors.New("broker: key not found")
	}
	cmd(f)
}

func (b *Broker) OnEvent(key string, h Handler) error {
	_, ok := b.direct[key]
	if !ok {
		b.direct[key] = make([]Handler, 0)
	}
	b.direct[key] = append(b.direct[key], h)
	return nil
}

func (b *Broker) Emit(key string) error {
	handlers, ok := b.direct[e.Key]
	if !ok {
		return errors.New("broker: key not found")
	}
	for _, handle := range handlers {
		handle()
	}
	return nil
}

