package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"

	"github.com/emersion/go-msgauth/dkim"
)

func lookupMX(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("error lookin up MX records: %v", err)
	}
	return mxRecords, nil
}

func sendMail(from, to string, data []byte) error {
	domain := strings.Split(to, "@")[1]

	mxRecords, err := lookupMX(domain)
	if err != nil {
		return err
	}

	for _, mx := range mxRecords {
		host := mx.Host

		for _, port := range []int{25, 587, 465} {
			address := fmt.Sprintf("%s:%d", host, port)

			var (
				client *smtp.Client
				err    error
				b      bytes.Buffer
			)

			switch port {
			case 465:
				// SMTPS
				tlsConfig := &tls.Config{ServerName: host}
				conn, err := tls.Dial("tcp", address, tlsConfig)
				if err != nil {
					continue
				}

				client, err = smtp.NewClient(conn, host)

			case 25, 587:
				// SMTP or SMTP with STARTTLS
				client, err = smtp.Dial(address)
				if err != nil {
					continue
				}
				if port == 587 {
					if err = client.StartTLS(&tls.Config{ServerName: host}); err != nil {
						client.Close()
						continue
					}
				}
			}

			if err := dkim.Sign(&b, bytes.NewReader(data), dkimOptions); err != nil {
				return fmt.Errorf("failed to sign email with DKIM: %v", err)
			}

			signedData := b.Bytes()

			// SMTP Conversation
			if err = client.Mail(from); err != nil {
				client.Close()
				continue
			}

			if err = client.Rcpt(to); err != nil {
				client.Close()
				continue
			}

			write, err := client.Data()
			if err != nil {
				client.Close()
				continue
			}

			_, err = write.Write(signedData)
			if err != nil {
				client.Close()
				continue
			}

			err = write.Close()
			if err != nil {
				client.Close()
				continue
			}

			client.Quit()

			return nil
		}
	}

	return fmt.Errorf("failed to send email to %s", to)
}

// load the DKIM private key
var dkimPrivateKey *rsa.PrivateKey

func init() {
	// load DKIM private key from a file
	privateKeyPEM, err := os.ReadFile("./private_key")
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		log.Fatalf("Failed to parse PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	dkimPrivateKey = privateKey
}

// DKIM Options
var dkimOptions = &dkim.SignOptions{
	Domain:   "example.com",
	Selector: "default",
	Signer:   dkimPrivateKey,
}
