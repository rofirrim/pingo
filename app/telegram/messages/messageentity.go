package messages

type MessageEntity struct {
	Type   string  `json:"type"`
	Offset int     `json:"offset"`
	Length int     `json:"length"`
	Url    *string `json:"url"`
	User   *User   `json:"user"`
}
