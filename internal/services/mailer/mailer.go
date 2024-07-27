package mailer

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Env interface {
	GetSMTPHost() string
	GetSMTPPort() string
	GetSMTPUser() string
	GetSMTPPassword() string
}

type Service struct {
	isEnabled bool
	host      string
	port      string
	user      string
	password  string
}

type mailerAuth struct {
	username string
	password string
}

func Init(env Env) (*Service, error) {
	host := strings.TrimSpace(env.GetSMTPHost())
	port := strings.TrimSpace(env.GetSMTPPort())
	user := strings.TrimSpace(env.GetSMTPUser())
	password := strings.TrimSpace(env.GetSMTPPassword())

	if host == "" || port == "" || user == "" || password == "" {
		return &Service{}, nil
	}

	return &Service{
		isEnabled: true,
		host:      host,
		port:      port,
		user:      user,
		password:  password,
	}, nil
}

func (m Service) Send(recipient string, subject string, text string) error {
	if !m.isEnabled {
		return nil
	}

	msg := []byte("To: " + recipient + "\r\n" +
		"From: " + m.user + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"" + "\n\r\n" +
		text + "\r\n")

	to := []string{recipient}
	auth := m.GetMailerAuth(m.user, m.password)
	return smtp.SendMail(m.host+":"+m.port, auth, m.user, to, msg)
}

func (m Service) GetMailerAuth(username, password string) smtp.Auth {
	return &mailerAuth{username, password}
}

func (a *mailerAuth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *mailerAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)

	if more {
		if command == "username" {
			return []byte(a.username), nil
		}
		if command == "password" {
			return []byte(a.password), nil
		}

		return nil, fmt.Errorf("unexpected server challenge: %s", command)
	}

	return nil, nil
}
