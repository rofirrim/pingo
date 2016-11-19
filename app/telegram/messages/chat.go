package messages

type Chat struct {
	Id                          int     `json:"id"`
	Type                        string  `json:"type"`
	Title                       *string `json:"title"`
	Username                    *string `json:"username"`
	FirstName                   *string `json:"first_name"`
	LastName                    *string `json:"last_name"`
	AllMembersAreAdministrators bool    `json:"all_members_are_administrators"`
}
