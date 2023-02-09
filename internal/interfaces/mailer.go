package interfaces

type Mailer interface {
	Send(recipient string, subject string, text string) error
}
