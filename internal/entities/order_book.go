package entities

type OrderBook struct {
	Id          string
	Amount      int
	PriceForOne float64
	BookId      string
	OrderId     string
	Book        Book `gorm:"foreignKey:BookId;references:Id"`
}
