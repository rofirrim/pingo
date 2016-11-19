package messages

type Update struct {
	UpdateId int               `json:"update_id"`
	Message  *Message `json:"message"`
}
