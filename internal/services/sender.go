package services

import "github.com/WildEgor/gNotifier/internal/domain"

type ISender interface {
	Send(req *domain.PushNotification) error
}

type Sender struct {
	sender ISender
}

func InitSender(s ISender) *Sender {
	return &Sender{
		sender: s,
	}
}

func (s *Sender) SetTransport(t ISender) *Sender {
	s.sender = t
	return s
}

func (s *Sender) Send(req *domain.PushNotification) error {
	err := s.sender.Send(req)
	return err
}
