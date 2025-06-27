package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomIcon struct {
	MyContent []byte
	MyName string
}
func (c *CustomIcon) Content() []byte {
	return c.MyContent
}
func (c *CustomIcon) Name() string {
	return c.MyName
}

func ChangeSVGIconFG(svg []byte, c color.RGBA) []byte {
	data := []byte{}
	red := fmt.Sprintf("%02x", c.R)
	green := fmt.Sprintf("%02x", c.G)
	blue := fmt.Sprintf("%02x", c.B)
	valueStart := -1
	for i, c := range svg {
		if c == '"' && valueStart == -1 {
			valueStart = i+1
			data = append(data, c)
		} else if c == '"' && valueStart > -1 {
			value := svg[valueStart:i]
			if string(value) == "#f3f3f3" {
				data = append(data, '#')
				data = append(data, red[0])
				data = append(data, red[1])
				data = append(data, green[0])
				data = append(data, green[1])
				data = append(data, blue[0])
				data = append(data, blue[1])
			} else {
				for _, c := range value {
					data = append(data, c)
				}
			}
			valueStart = -1
		}
		if valueStart == -1 {
			data = append(data, c)
		}
	}
	return data
}

func NewDeleteIcon() fyne.Resource {
	icon := theme.DeleteIcon()
	data := icon.Content()
	data = ChangeSVGIconFG(data, color.RGBA{0xff, 0x00, 0x00, 0x00})
	myIcon := &CustomIcon{
		MyName: "MyDeleteIcon",
		MyContent: data,
	}
	return myIcon
}

func NewEditIcon(color color.RGBA, name string) fyne.Resource {
	data := theme.DocumentCreateIcon().Content()
	data = ChangeSVGIconFG(data, color)
	myIcon := &CustomIcon{
		MyName: name,
		MyContent: data,
	}
	return myIcon
}
