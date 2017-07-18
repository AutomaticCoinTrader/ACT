package utility

import (
	"github.com/pkg/errors"
	"net"
	"crypto/tls"
	"net/mail"
	"net/smtp"
	"strings"
	"fmt"
	"log"
)

const (
	SMTPAuthUnkown  SMTPAuthType = ""
	SMTPAuthPLAIN 	SMTPAuthType = "PLAIN"
	AMTPAuthCRAMMD5 SMTPAuthType = "CRAM-MD5"
)

type SMTPAuthType string

func (s SMTPAuthType) String() (string) {
	return string(s)
}

func GetSMTPAuthType(authType string) (SMTPAuthType) {
	switch strings.ToUpper(authType) {
	case SMTPAuthPLAIN.String():
		return SMTPAuthPLAIN
	case AMTPAuthCRAMMD5.String():
		return AMTPAuthCRAMMD5
	default:
		return SMTPAuthUnkown
	}
}

// SMTPClient is smtp clinet
type SMTPClient struct {
	hostPort    string
	username    string
	password    string
	authType    SMTPAuthType
	useTLS      bool
	useStartTLS bool
	from        string
    to          string
}

func (s *SMTPClient) SendMail(subject string, body string) (error) {
	from := mail.Address{
		Address: s.from,
	}
	toList, err := mail.ParseAddressList(s.to)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not parse mail address list (to = %v)", s.to))
	}
	message := ""
	message += fmt.Sprintf("From: %s\r\n", s.from)
	message += fmt.Sprintf("To: %s\r\n", s.to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "\r\n" + body

	host, _, _ := net.SplitHostPort(s.hostPort)

	var auth smtp.Auth = nil
	if s.username != "" {
		if s.authType == SMTPAuthPLAIN {
			auth = smtp.PlainAuth("", s.username, s.password, host)
		} else if s.authType == AMTPAuthCRAMMD5 {
			auth = smtp.CRAMMD5Auth(s.username, s.password)
		}
	}

	var conn net.Conn
	if s.useTLS {
		tlsContext := &tls.Config {
			ServerName: host,
			InsecureSkipVerify: false,
		}
		conn, err = tls.Dial("tcp", s.hostPort, tlsContext)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not connect mail host with tls (host port = %v, use tls = %v)", s.hostPort, s.useTLS))
		}
	} else {
		conn, err = net.Dial("tcp", s.hostPort)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not connect mail host (host port = %v)", s.hostPort))
		}
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not create smtp client (host port = %v)", s.hostPort))
	}

	if s.useStartTLS {
		tlsconfig := &tls.Config {
			ServerName: host,
			InsecureSkipVerify: false,
		}
		if err := client.StartTLS(tlsconfig); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not start tls (host port = %v, use start tls = %v)", s.hostPort, s.useStartTLS))
		}
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can not authentication (host port = %v authType = %v, username = %v password = %v)", s.hostPort, s.authType, s.username, s.password))
		}
	}

	if err = client.Mail(from.Address); err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not send MAIL command (from = %v)", from.Address))
	}

	var emails []string
	for _,  to := range toList {
		emails = append(emails, to.Address)
	}
	recept := strings.Join(emails, ",")
	if err = client.Rcpt(recept); err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not send RCPT command (recept = %v)", recept))
	}

	w, err := client.Data()
	if err != nil {
		return errors.Wrap(err, "can not send DATA command")
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		w.Close()
		return errors.Wrap(err, "can not write message")
	}

	err = w.Close()
	if err != nil {
		log.Printf("can not close message writer (reason = %v)", err)
	}

	err = client.Quit()
	if err != nil {
		log.Printf("can not send QUIT command (reason = %v)", err)
	}

	return nil
}

// NewSMTPClient is create smtp client
func NewSMTPClient(hostPort string, username string, password string, authtype SMTPAuthType, useTLS bool, useStartTLS bool, from string, to string) (n *SMTPClient) {
	return &SMTPClient{
		hostPort: hostPort,
		username: username,
		password: password,
		authType: authtype,
		useTLS: useTLS,
		useStartTLS: useStartTLS,
		from: from,
		to: to,
	}
}
