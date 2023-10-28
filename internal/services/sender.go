package services

import "github.com/WildEgor/gNotifier/internal/domain"

// ISender Interface for future Senders
type ISender interface {
	Send(req *domain.PushNotification) error
}

// Sender Strategy Pattern
type Sender struct {
	sender ISender
}

// NewSender Constructor
func NewSender(s ISender) *Sender {
	return &Sender{
		sender: s,
	}
}

// SetTransport Allow change transport
func (s *Sender) SetTransport(t ISender) *Sender {
	s.sender = t
	return s
}

// Send Execute send method
func (s *Sender) Send(req *domain.PushNotification) error {
	err := s.sender.Send(req)
	return err
}
