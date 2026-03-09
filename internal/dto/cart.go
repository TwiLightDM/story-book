package dto

type CartRequest struct {
	BookId string `json:"book_id"`
	Amount int    `json:"amount"`
}
