package main // import "github.com/psanford/unsubscribe"

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
)

var (
	smtpAddress  = flag.String("stmp-address", "localhost:25", "STMP server address")
	useTLS       = flag.Bool("tls", true, "Use STARTTLS")
	authUser     = flag.String("auth-user", "", "Authenticate to server as user")
	authPassword = flag.String("auth-password", "", "Authenticate to server with password")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatalf("usage: %s <from_address> <to_address>\n", os.Args[0])
	}

	from := args[0]
	to := args[1]

	headers := []Header{
		{"From", from},
		{"To", to},
		{"Subject", "unsubscribe"},
		{"MIME-Version", "1.0"},
		{"Content-Type", "text/plain; charset=\"utf-8\""},
		{"Content-Transfer-Encoding", "7bit"},
	}

	message := make([]byte, 0)
	for _, h := range headers {
		message = append(message, []byte(h.String())...)
	}
	message = append(message, []byte("\r\n")...)
	message = append(message, []byte("unsubscribe")...)

	c, err := smtp.Dial(*smtpAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	host, _, err := net.SplitHostPort(*smtpAddress)
	if err != nil {
		log.Fatal(err)
	}

	tc := tls.Config{
		ServerName: host,
	}

	if err = c.StartTLS(&tc); err != nil {
		log.Fatal(err)
	}

	if *authUser != "" {
		a := smtp.PlainAuth("", *authUser, *authPassword, host)
		if err = c.Auth(a); err != nil {
			log.Fatal(err)
		}
	}

	if err = c.Mail(from); err != nil {
		log.Fatal(err)
	}

	if err = c.Rcpt(to); err != nil {
		log.Fatal(err)
	}

	w, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	if _, err = w.Write([]byte(message)); err != nil {
		log.Fatal(err)
	}
	if err = w.Close(); err != nil {
		log.Fatal(err)
	}

	c.Quit()
}

type Header struct {
	Name  string
	Value string
}

func (h Header) String() string {
	return fmt.Sprintf("%s: %s\r\n", h.Name, h.Value)
}
