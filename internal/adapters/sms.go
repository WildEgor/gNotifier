package adapters

type ISMSAdapter interface {
	Send() error
}

type SMSAdapter struct {
}

func NewSMSAdapter() *SMSAdapter {
	return &SMSAdapter{}
}

func (s *SMSAdapter) Send() error {
	return nil
}
