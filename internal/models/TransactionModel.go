package models

// convert to mysql later.
type TransactionModel struct {
	// change to id later
	Id            string  `json:"id"`
	Organization  string  `json:"organization"`
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	CreatedAt     string  `json:"created_at"`
	Gateway       string  `json:"gateway"`
	Status        string  `json:"status"`
	CallbackURLS  []string  `json:"callback_urls,omitempty"`
}
