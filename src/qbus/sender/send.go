package send

import (
	"github.com/qbusapi/app/models"
)

type Sender interface {
	Send(message *models.MessageModel) (err error)
}

type AllSend struct {
	SMSSender   Sender
	EmailSender Sender
}

func NewAllSend(smsSender, emailSender Sender) Sender {
	return &AllSend{
		SMSSender:   smsSender,
		EmailSender: emailSender,
	}
}

func (s *AllSend) Send(message *models.MessageModel) (err error) {
	switch message.Type {
	case models.MessageTypeMail:
		return s.EmailSender.Send(message)
	case models.MessageTypeSMS:
		return s.SMSSender.Send(message)
	default:
		return UnsupportMessageTypeError
	}
}
