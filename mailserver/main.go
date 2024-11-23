package main

import (
	"fmt"
	"io"
	"log"
	"time"

	smtp "github.com/emersion/go-smtp"
)

// type Backend implements SMTP methods
type Backend struct{}

// type Session returns a Session after initiation with EHLO
type Session struct {
	From string
	To   []string
}

func main() {
	srv := smtp.NewServer(&Backend{})

	srv.Addr = "localhost:2525"
	srv.Domain = "localhost"
	srv.ReadTimeout = 10 * time.Second
	srv.WriteTimeout = 10 * time.Second
	srv.MaxMessageBytes = 1024 * 1024
	srv.MaxRecipients = 50
	srv.AllowInsecureAuth = true

	// start the SMTP server
	log.Println("Starting server at", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

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
	if msg, err := io.ReadAll(r); err != nil {
		return err
	} else {
		fmt.Println("Received Message: ", string(msg))
		for _, recipient := range s.To {
			if err := sendMail(s.From, recipient, msg); err != nil {
				fmt.Printf("Failed to send email to %s: %v", recipient, err)
			} else {
				fmt.Printf("Email sent succesfully to %s", recipient)
			}
		}
		return nil
	}
}

func (s *Session) AuthPlain(username, password string) error {
	if username != "testuser" || password != "testpass" {
		return fmt.Errorf("invalid username or password")
	}
	return nil
}

func (s *Session) Reset() {
	s.From = ""
	s.To = []string{}
}

func (s *Session) Logout() error {
	return nil
}
