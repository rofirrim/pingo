package messages

type User struct {
	Id        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Username  *string `json:"username"`
}
