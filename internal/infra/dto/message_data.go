package dto

type MessageData struct {
	QRCode string `json:"qr_code,omitempty"`
	ID     int    `json:"id,omitempty"`
}

func NewMessageData(qrCode string, id int) *MessageData {
	return &MessageData{
		QRCode: qrCode, //qr code
		ID:     id,     //id do evento
		//public id
		//status - confirmado
	}
}
