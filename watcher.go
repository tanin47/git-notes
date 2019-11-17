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
	checkInterval time.Duration
	delayBeforeFiringEvent time.Duration
	delayAfterFiringEvent time.Duration
}

func (f *GitWatcher) Stop() {
	f.running = false
}

func (f *GitWatcher) Check(path string, channel chan string) {
	dirty, err := f.git.IsDirty(path)

	if err != nil {
		log.Printf("Failed to get state. Error: %v", err)
	}

	if dirty {
		log.Printf("Changes have been detected.")
		time.Sleep(f.delayBeforeFiringEvent)
		channel <- path
		time.Sleep(f.delayAfterFiringEvent)
	}
}

func (f *GitWatcher) Watch(path string, channel chan string) {
	f.running = true
	go func() {
		for f.running {
			time.Sleep(f.checkInterval)
			f.Check(path, channel)
		}
	}()

}
