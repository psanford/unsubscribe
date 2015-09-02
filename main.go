package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: %s <from_address> <to_address>\n", os.Args[0])
	}

	from := os.Args[1]
	to := os.Args[2]

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

	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
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
