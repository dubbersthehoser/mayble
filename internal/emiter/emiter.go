package emiter


type Emiter struct {
	event map[string][]func(any)
	request map[string]func(func(any))
}

func NewEmiter() *Emiter {
	return &Emiter{
		event: make(map[string[]Handler),
		request: make(map[string]func(func(any))),
	}
}

func (e *Emiter) OnRequest(key string, f func(func(any))) {
	e.request[key] = f
}

func (e *Emiter) Request(key string, f func(any)) error {
	cmd, ok := e.request[key]
	if !ok {
		return errors.New("emiter: key not found")
	}
	cmd(f)
}

func (e *Emiter) OnEvent(key string, handler func(any)) {
	_, ok := e.event[key]
	if !ok {
		e.event[key] = make([]Handler, 0)
	}
	e.event[key] = append(e.evetn[key], h)
}

func (e *Emiter) Emit(key string, data any) error {
	handlers, ok := e.event[e.Key]
	if !ok {
		return errors.New("emiter: key not found")
	}
	for _, handle := range handlers {
		handle(key, data)
	}
	return nil
}

