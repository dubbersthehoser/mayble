package model


import (
	"time"

)


type Model struct {
}

type  Book struct {
	Title string
	Author string
	Genre string
	Ratting int
}

type Loan struct {
	Name string
	Date time.Time
}


func GetRattingStrings() []string {
	return []string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}
}

