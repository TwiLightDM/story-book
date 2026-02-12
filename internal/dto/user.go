package dto

type UserRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type UserResponse struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Surname  string `json:"surname,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Role     string `json:"role,omitempty"`
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
	Points   int    `json:"points,omitempty"`
	Error    string `json:"error,omitempty"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users,omitempty"`
	Error string         `json:"error,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Error        string `json:"error,omitempty"`
}

type SignUpResponse struct {
	User         UserResponse `json:"user,omitempty"`
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	Error        string       `json:"error,omitempty"`
}
