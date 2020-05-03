package users

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
