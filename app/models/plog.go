package models

import "time"

type Plog struct {
    Id           int
	Text         string
	Protagonista User
	Autor        User
	Titol        string
	Data         time.Time
    Nota         float32
}
