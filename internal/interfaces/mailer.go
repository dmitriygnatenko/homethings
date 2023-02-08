package interfaces

type IMailer interface {
	Send(recipient string, subject string, text string) error
}
