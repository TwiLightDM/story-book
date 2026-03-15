package dto

type UserRequest struct {
	Name     string  `json:"name"`
	Surname  string  `json:"surname"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Password string  `json:"password"`
	Question string  `json:"question"`
	Answer   string  `json:"answer"`
	Address  *string `json:"address"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Question string `json:"question"`
	Points   int    `json:"points"`
	Address  string `json:"address,omitempty"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}
