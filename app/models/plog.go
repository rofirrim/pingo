package models

type Plog struct {
    Id           int
	Text         string
	Protagonista User
	Autor        User
	Titol        string
	Dia          string
	Hora         string
    Nota         float32
}
