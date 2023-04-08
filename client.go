package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/emersion/go-imap/client"
)

/*
creates a client and logins
*/
func NewClient(config Mailbox) (*client.Client, error) {
	log.Println("Connecting to server...")
	// Connect to server
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", config.Host, config.Port), &tls.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Connected")

	// c.Timeout = time.Second * 10
	// Login
	if err = c.Login(config.Username, config.Password); err != nil {
		return nil, err
	}

	log.Println("Logged in")

	return c, nil
}
