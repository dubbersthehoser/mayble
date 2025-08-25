package config

import (
	"bytes"
	"fmt"
	//"os"
	"testing"
)

func compareSections(e, a Sections) error {
	eSections := e
	aSections := a
	for eSection, eProps := range eSections {
		aProps, ok := aSections[eSection]
		if !ok {
			return fmt.Errorf("section was not found: '%s'", eSection)
			continue
		}

		for k, v := range eProps {
			value, ok := aProps[k]
			if !ok {
				return fmt.Errorf("prop key was not found: '%s'", k)
			} else if v != value{
				return fmt.Errorf("expect value '%s', got '%s'", v, value)
			}

		}
	}
	return nil
	
}

func TestReadWrite(t *testing.T) {

	
	jdata := []byte(`{
	"section_0": {"key_0": "value_0"},
	"section_1": {"key_0": "value_0", "key_1": "value_1", "key_2": "value_2"}
	}`)
	sections := Sections{
		"section_0": Properties {
			"key_0": "value_0",
		},
		"section_1": Properties {
			"key_0": "value_0",
			"key_1": "value_1",
			"key_2": "value_2",
		},
	}

	ReadFailed := true

	t.Run("Read", func(t *testing.T) {
		var (
			expect Sections = sections
		)
		r := bytes.NewReader(jdata)
		aSections, err := Read(r)
		if err != nil {
			t.Fatal(err)
		}
		if err := compareSections(expect, aSections); err != nil {
			t.Error(err)
		}
		ReadFailed = t.Failed()
	})
	t.Run("Write", func(t *testing.T){
		if ReadFailed {
			t.Fatal("Read() Failed. Unable to do Write() testing.")
		}
		data := []byte{}
		buf := bytes.NewBuffer(data)
		if err := Write(sections, buf); err != nil {
			t.Fatal(err)
		}

		aSections, err := Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if err := compareSections(sections, aSections); err != nil {
			t.Error(err)
		}
	})
	
}





