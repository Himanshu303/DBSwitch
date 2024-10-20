package models

import "time"

type Student struct {
	ID    int        `json:"_id"`
	Name  string     `json:"name"`
	Bdate *time.Time `json:"bdate"`
	Marks int        `json:"marks"`
	Gpa   float32    `json:"gpa"`
}
