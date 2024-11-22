package enums

type TransactionStatus string

const (
	PENDING    TransactionStatus = "PENDING"
	SETTLEMENT TransactionStatus = "SETTLEMENT"
	REJECTED   TransactionStatus = "REJECTED"
)
