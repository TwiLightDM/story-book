package dto

type OrderRequest struct {
	Address *string `json:"address"`
	ShopId  *string `json:"shop_id"`
	Code    *string `json:"code"`
}

type OrderResponse struct {
	Id           string               `json:"id"`
	Status       string               `json:"status"`
	DeliveryType string               `json:"delivery_type"`
	Address      string               `json:"address,omitempty"`
	Cost         float64              `json:"cost"`
	Points       int                  `json:"points"`
	CreatedAt    string               `json:"created_at"`
	Shop         ShopResponse         `json:"shop,omitempty"`
	OrderBooks   []OrderBooksResponse `json:"order_books"`
}

type OrderBooksResponse struct {
	Amount      int          `json:"amount"`
	PriceForOne float64      `json:"price_for_one"`
	Book        BookResponse `json:"book"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
}

type StatusRequest struct {
	Status string `json:"status"`
}
