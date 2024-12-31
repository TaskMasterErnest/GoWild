package main

import (
	"fmt"
	"io"

	"github.com/emersion/go-smtp"
)

// implement SMTP methods
type Backend struct{}

// implement Sessions
type Session struct {
	From string
	To   []string
}

// start a new session with a new connection
func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	fmt.Println("Mail From: ", from)
	s.From = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	fmt.Println("Recipient To: ", to)
	s.To = append(s.To, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if message, err := io.ReadAll(r); err != nil {
		return err
	} else {
		fmt.Println("Received Message: ", string(message))
	}
	return nil
}

// adding authentication mechanisms
func (s *Session) AuthMechanisms() []string {
	return []string{sasl.}
}
