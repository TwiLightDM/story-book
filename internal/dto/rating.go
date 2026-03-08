package dto

type RatingRequest struct {
	Stars  int    `json:"stars"`
	BookId string `json:"book_id"`
}
