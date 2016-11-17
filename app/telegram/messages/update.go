package messages

type Update struct {
	UpdateId int
	Message  *messages.Message
	// EditedMessage *messages.Message
	// InlineQuery *messages.InlineQuery
	// ChosenInlineResult *messages.ChosenInlineResult
	// CallbackQuery *messages.CallbackQuery
}
