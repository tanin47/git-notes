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
			state, err := f.git.GetState(path)

			if err != nil {
				log.Printf("Failed to get state. Error: %v", err)
			}

			if state == Dirty {
				log.Printf("Changes have been detected.")
				time.Sleep(3 * time.Second)
				channel <- path
			}
		}
	}()

}
