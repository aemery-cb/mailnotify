package main

import (
	"context"
	"log"
	"time"

	"github.com/emersion/go-imap/client"
)

type MailboxWatcher struct {
	c        *client.Client
	mailbox  string
	postHook string
}

/*
logs into the mailbox and establishes a
new watcher with client.
*/
func NewMailWatcher(c *client.Client, mailbox string, postHook string) *MailboxWatcher {
	return &MailboxWatcher{c: c, mailbox: mailbox, postHook: postHook}
}

func (w *MailboxWatcher) Logout() error {
	log.Println("logging out")
	w.c.Timeout = time.Second * 10
	// TODO: logging out is broken and blocking.
	return w.c.Logout()
}

func (w *MailboxWatcher) Watch(ctx context.Context) error {

	_, err := w.c.Select(w.mailbox, false)
	if err != nil {
		return err
	}

	updates := make(chan client.Update)
	w.c.Updates = updates
	done := make(chan error, 1)
	go func() {
		done <- w.c.Idle(ctx.Done(), nil)
	}()
	log.Println("Waiting for updates")

	for {
		select {
		case update := <-updates:
			log.Println("New update:", update)
			w.HandleMessageType(update)
		case err := <-done:
			if err != nil {
				return err
			}
			log.Println("Not idling anymore")
			return nil
		}
	}
}

/*
executes the appropriate action
*/
func (w *MailboxWatcher) HandleMessageType(update client.Update) {
	// switch message := update.(type) {
	// case *client.MailboxUpdate:
	// 	log.Println("mailbox update")
	// case *client.ExpungeUpdate:
	// 	log.Println("expunge update")
	// case *client.MessageUpdate:
	// 	log.Println("message update")
	// case *client.StatusUpdate:
	// 	log.Println("status update")
	// }

	// assume everything just wants the postHook to run.
	err := ExecuteHook(w.postHook)
	if err != nil {
		log.Println("failed to execute")
		log.Println(err)
	} else {
		log.Println("successfully exec hook")
	}

}
