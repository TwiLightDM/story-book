package dto

type BookRequest struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Year        int     `json:"year"`
	Cost        float64 `json:"cost"`
	Discount    int     `json:"discount"`
	Publisher   string  `json:"publisher"`
	Description string  `json:"description"`
	Amount      int     `json:"amount"`
	Image       string  `json:"image"`
}

type BookResponse struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Year        int     `json:"year"`
	Cost        float64 `json:"cost"`
	Discount    int     `json:"discount,omitempty"`
	Publisher   string  `json:"publisher"`
	Description string  `json:"description,omitempty"`
	Amount      int     `json:"amount"`
	Image       string  `json:"image,omitempty"`
}

type BookListResponse struct {
	Books []BookResponse `json:"books"`
}
