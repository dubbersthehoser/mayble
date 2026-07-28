package viewmodel

const (
	BodyNoData int = iota
	BodyTable
	BodyBookEdit
	BodyBookCreate
	BodyInfo
	BodyManual
)

type Body struct {
	prev  int
	value int
	l     []func()

	hideHandlers map[int]func()
	showHandlers map[int]func()
}

func (b *Body) RegisterHandlers(body int, hide, show func()) {
	if b.hideHandlers == nil {
		b.hideHandlers = make(map[int]func())
	}
	if b.showHandlers == nil {
		b.showHandlers = make(map[int]func())
	}

	b.hideHandlers[body] = hide
	b.showHandlers[body] = show
}

func (b *Body) Value() int {
	return b.value
}

func (b *Body) Set(v int) {
	for _, h := range b.hideHandlers {
		h()
	}
	if b.showHandlers != nil {
		b.showHandlers[v]()
	}
	b.prev = b.value
	b.value = v
	b.notify()
}

func (b *Body) Back() {
	b.Set(b.prev)
}

func (b *Body) AddListener(fn func()) {
	if b.l == nil {
		b.l = make([]func(), 0)
	}
	b.l = append(b.l, fn)
}

func (b *Body) notify() {
	for _, fn := range b.l {
		fn()
	}
}
