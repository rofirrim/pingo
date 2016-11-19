package messages

type Message struct {
	MessageId       int              `json:"id"`
	From            *User            `json:"from"`
	Date            int              `json:"date"`
	Chat            Chat             `json:"chat"`
	ForwardFrom     *User            `json:"forward_from"`
	ForwardFromChat *Chat            `json:"forward_from_chat"`
	ForwardDate     *int             `json:"forward_date"`
	Text            *string          `json:"text"`
	Entities        *[]MessageEntity `json:"entities"`
}
