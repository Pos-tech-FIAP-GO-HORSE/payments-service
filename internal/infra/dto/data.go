package dto

type ResponseCreatePayment struct {
	QRCode string `json:"qr_code"`
	ID     int    `json:"id"`
}

type ResponseStatusPayment struct {
	ID                int     `json:"id"`
	Status            string  `json:"status"`
	StatusDetail      string  `json:"status_detail"`
	TransactionAmount float64 `json:"transaction_amount"`
	PaymentMethodId   string  `json:"payment_method_id"`
}
