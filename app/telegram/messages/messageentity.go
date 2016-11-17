package messages

type MessageEntity struct {
	Type   string
	Offset int
	Length int
	Url    *string
	User   *messages.User
}
