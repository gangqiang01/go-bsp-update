package v1

type Auth struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}
