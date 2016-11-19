package commands

type SendMessage struct {
	Method    string `json:"method"`
	ChatId    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func NewSendMessage() SendMessage {
	var res SendMessage
	res.Method = "sendMessage"
	return res
}
