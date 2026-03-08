package dto

type GenreRequest struct {
	Genre  string `json:"genre"`
	BookId string `json:"book_id"`
}

type GenreResponse struct {
	Genre string `json:"genre"`
}
