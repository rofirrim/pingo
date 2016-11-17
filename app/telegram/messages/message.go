package messages

type Message struct {
    MessageId	int	//Unique message identifier
    From	*messages.User	//. Sender, can be empty for messages sent to channels
    Date	int	//Date the message was sent in Unix time
    Chat	messages.Chat	// Conversation the message belongs to
    ForwardFrom	*messages.User	//For forwarded messages, sender of the original message
    ForwardFromChat	*messages.Chat	//For messages forwarded from a channel, information about the original channel
    ForwardDate	*int	//. For forwarded messages, date the original message was sent in Unix time
    ReplyToMessage	*messages.Message	//. For replies, the original message. Note that the Message object in this field will not contain further replyToMessage fields even if it itself is a reply.
    EditDate	*int	//. Date the message was last edited in Unix time
    Text	*string	//. For text messages, the actual UTF-8 text of the message, 0-4096 characters.
    Entities	*[]messages.MessageEntity	//. For text messages, special entities like usernames, URLs, bot commands, etc. that appear in the text
    // Audio	*messages.Audio	//. Message is an audio file, information about the file
    // Document	*messages.Document	//. Message is a general file, information about the file
    // Game	*messages.Game	//. Message is a game, information about the game. More about games »
    // Photo	*[]messages.PhotoSize	//. Message is a photo, available sizes of the photo
    // Sticker	*messages.Sticker	//. Message is a sticker, information about the sticker
    // Video	*messages.Video	//. Message is a video, information about the video
    // Voice	*messages.Voice	//. Message is a voice message, information about the file
    // Caption	*string	//. Caption for the document, photo or video, 0-200 characters
    // Contact	*messages.Contact	//. Message is a shared contact, information about the contact
    // Location	*messages.Location	//. Message is a shared location, information about the location
    // Venue	*messages.Venue	//. Message is a venue, information about the venue
    // NewChatMember	*messages.User	//. A new member was added to the group, information about them (this member may be the bot itself)
    // LeftChatMember	*messages.User	//. A member was removed from the group, information about them (this member may be the bot itself)
    // NewChatTitle	*string	//. A chat title was changed to this value
    // NewChatPhoto	*[]messages.PhotoSize	//. A chat photo was change to this value
    // DeleteChatPhoto	*bool	//. Service message: the chat photo was deleted
    // GroupChatCreated	*bool	//. Service message: the group has been created
    // SupergroupChatCreated	*bool	//. Service message: the supergroup has been created. This field can‘t be received in a message coming through updates, because bot can’t be a member of a supergroup when it is created. It can only be found in replyToMessage if someone replies to a very first message in a directly created supergroup.
    // ChannelChatCreated	*bool	//. Service message: the channel has been created. This field can‘t be received in a message coming through updates, because bot can’t be a member of a channel when it is created. It can only be found in replyToMessage if someone replies to a very first message in a channel.
    // MigrateToChatId	*int	//. The group has been migrated to a supergroup with the specified identifier. This number may be greater than 32 bits and some programming languages may have difficulty/silent defects in interpreting it. But it smaller than 52 bits, so a signed 64 bit integer or double-precision float type are safe for storing this identifier.
    // MigrateFromChatId	*int	//. The supergroup has been migrated from a group with the specified identifier. This number may be greater than 32 bits and some programming languages may have difficulty/silent defects in interpreting it. But it smaller than 52 bits, so a signed 64 bit integer or double-precision float type are safe for storing this identifier.
    // PinnedMessage	*messages.Message	//. Specified message was pinned. Note that the Message object in this field will not contain further replyToMessage fields even if it is itself a reply.
}
