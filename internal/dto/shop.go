package dto

type ShopRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type ShopResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type ShopListResponse struct {
	Shops []ShopResponse `json:"shops"`
}
