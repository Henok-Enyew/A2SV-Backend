package entity

type User struct {
	ID       string
	Username string
	Password string
	Role     string
}

func NewUser(id, username, password, role string) *User {
	return &User{
		ID:       id,
		Username: username,
		Password: password,
		Role:     role,
	}
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u *User) PromoteToAdmin() {
	u.Role = "admin"
}


