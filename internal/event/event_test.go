package event

import (
	"testing"
)


func TestEvents(t *testing.T) {

	const (
		PRINT string = "PRINT"
		INC          = "INC"
		DEC          = "DEC"
		SET          = "SET"
	)
	
	count := 0
	emiter := NewEventEmiter()

	emiter.On(
		PRINT, 
		func(data any) {
			t.Log("Print Event Ran")
		},
	)

	emiter.On(
		INC,
		func(data any) {
			count += 1
		},
	)

	emiter.On(
		DEC,
		func(data any) {
			count -= 1
		},
	)

	emiter.On(
		SET,
		func(data any) {
			i, ok := data.(int)
			if ok {
				count = i
			}
		},
	)

	emiter.Emit(PRINT, nil)

	emiter.Emit(INC, nil)
	emiter.Emit(INC, nil)
	emiter.Emit(INC, nil)

	if count != 3 {
		t.Errorf("expect 3, got %d", count)
	}
	emiter.Emit(DEC, nil)
	emiter.Emit(DEC, nil)
	emiter.Emit(DEC, nil)
	emiter.Emit(DEC, nil)

	if count != -1 {
		t.Errorf("expect -1, got %d", count)
	}

	emiter.Emit(SET, 10)

	if count != 10 {
		t.Errorf("expect 10, got %d", count)
	}
}
