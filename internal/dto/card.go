package dto

type CardRequest struct {
	NumberOfCard   string `json:"number_of_card"`
	ExpirationDate string `json:"expiration_date"`
	Cvv            string `json:"cvv"`
}

type CardResponse struct {
	NumberOfCard string `json:"number_of_card"`
}

type CardListResponse struct {
	Cards []CardResponse `json:"cards"`
}
