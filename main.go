package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	confPath := flag.String("conf", "./mailnotify.json", "Configuration file")
	flag.Parse()

	configs, err := NewConfig(confPath)
	if err != nil {
		log.Fatal(err)
	}

	watchers := make([]*MailboxWatcher, len(*configs))

	for i, config := range *configs {
		c, err := NewClient(config)
		if err != nil {
			log.Fatal(err)
		}
		// TODO: add support for watching multiple mailboxes here
		watchers[i] = NewMailWatcher(c, config.Box, config.PostHook)
	}

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	for _, watcher := range watchers {
		wg.Add(1)
		go func(watcher *MailboxWatcher, ctx context.Context) {
			if err := watcher.Watch(ctx); err != nil {
				log.Println(err)
			}
			if err := watcher.Logout(); err != nil {
				log.Println(err)
			}
			wg.Done()
		}(watcher, ctx)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	go func() {
		<-shutdown
		log.Println("gracefully shutting down")
		cancel()
		wg.Wait()

		os.Exit(1)
	}()

	shutdown <- os.Interrupt
	wg.Wait()
}
