package adapters

type ISMTPAdapter interface {
	Send() error
}

type SMTPAdapter struct {
}

func NewSMTPAdapter() *SMTPAdapter {
	return &SMTPAdapter{}
}

func (s *SMTPAdapter) Send() error {
	return nil
}
