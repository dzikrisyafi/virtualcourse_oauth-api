package users

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
