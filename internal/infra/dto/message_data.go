package dto

import "strconv"

type MessageData struct {
	ID       string `json:"id"`
	PublicID string `json:"public_id"`
	Status   string `json:"status"`
	QRCode   string `json:"qr_code"`
}

func NewMessageData(qrCode string, id int, publicID string) *MessageData {
	return &MessageData{
		QRCode:   qrCode,
		ID:       strconv.Itoa(id),
		PublicID: publicID,
		Status:   "Confirmado",
	}
}
