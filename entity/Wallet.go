package entity

type Wallet struct {
	Id         string  `json:"id"`
	CustomerId string  `json:"customer_id"`
	Balance    float64 `json:"balance"`
}
