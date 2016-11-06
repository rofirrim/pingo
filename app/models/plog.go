package models

import "time"

type Plog struct {
	Text         string
	Protagonista User
	Autor        User
	Titol        string
	Data         time.Time
}
