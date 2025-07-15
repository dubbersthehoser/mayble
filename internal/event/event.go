package event

type EventEmiter struct  {
	listeners map[string][]func(data any)
}
func NewEventEmiter() *EventEmiter {
	e := &EventEmiter{
		listeners: map[string][]func(data any){},
	}
	return e
}
func (e *EventEmiter) On(s string, fn func(data any)) {
	handlers, ok := e.listeners[s]
	if !ok {
		handlers = []func(data any){}
	}
	handlers = append(handlers, fn)
	e.listeners[s] = handlers
}
func (e *EventEmiter) Emit(s string, data any){
	for _, handler := range e.listeners[s] {
		handler(data)
	}
}
