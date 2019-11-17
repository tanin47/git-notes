package main

import (
	"log"
	"time"
)

type Watcher interface {
	Watch(path string, channel chan string)
}

type GitWatcher struct {
	git Git
	running bool
}

func (f *GitWatcher) Stop() {
	f.running = false
}

func (f *GitWatcher) Watch(path string, channel chan string) {
	f.running = true
	go func() {
		for {
			time.Sleep(10 * time.Second)
			dirty, err := f.git.IsDirty(path)

			if err != nil {
				log.Printf("Failed to get state. Error: %v", err)
			}

			if dirty {
				log.Printf("Changes have been detected.")
				time.Sleep(3 * time.Second)
				channel <- path
			}
		}
	}()

}
