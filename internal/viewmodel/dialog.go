package viewmodel


type ShowError struct {
	isOpen bool
	Err    error
	l []func()
}

func (se *ShowError) Show(err error) {
	se.Err = err
	se.notify()
}

func (se *ShowError) AddListener(h func()) {
	if se.l == nil {
		se.l = make([]func(), 0)
	}
	se.l = append(se.l, h)
}

func (se *ShowError) notify() {
	for _, h := range se.l {
		h()
	}
}

