package models

type Plog struct {
	Id           int
	RawText      string
	Text         string
	Protagonista User
	Autor        User
	RawTitol     string
	Titol        string
	Dia          string
	DiaYMD       string
	Hora         string
	Nota         float32

	AllowEdit bool
}

func (plog Plog) Description() string {
	return plog.Text
}

// Used for uploading
type PlogData struct {
	Text         string
	Protagonista int
	Autor        int
	Titol        string
	Data         int64 // Unix epoch
}
