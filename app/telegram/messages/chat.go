package messages

type Chat struct {
	Id                          int
	Type                        string
	Title                       *string
	Username                    *string
	FirstName                   *string
	LastName                    *string
	AllMembersAreAdministrators bool
}
