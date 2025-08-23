package event

type EventData interface {
	EventData()
}

type EventEmiter struct  {
	listeners map[string][]func(data EventData)
}
func NewEventEmiter() *EventEmiter {
	e := &EventEmiter{
		listeners: map[string][]func(data EventData){},
	}
	return e
}
func (e *EventEmiter) On(s string, fn func(data EventData)) {
	handlers, ok := e.listeners[s]
	if !ok {
		handlers = []func(data EventData){}
	}
	handlers = append(handlers, fn)
	e.listeners[s] = handlers
}
func (e *EventEmiter) Emit(s string, data EventData){
	for _, handler := range e.listeners[s] {
		handler(data)
	}
}
