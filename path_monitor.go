package main

import (
	"log"
	"time"
)

type PathMonitor interface {
	StartMonitoring(repoPath string, watcher Watcher, git Git)
	scheduleUpdate(repoPath string, channel chan string)
}

type GitRepoMonitor struct {
	scheduledUpdateInterval time.Duration
}

func (g *GitRepoMonitor) scheduleUpdate(repoPath string, channel chan string) {
	time.AfterFunc(g.scheduledUpdateInterval, func() {
		channel <- repoPath
		g.scheduleUpdate(repoPath, channel)
	})
}

func (g *GitRepoMonitor) StartMonitoring(repoPath string, watcher Watcher, git Git) {
	var channel = make(chan string)
	err := git.Sync(repoPath)
	if err != nil {
		log.Printf("Syncing failed. Err: %v", err)
	}
	g.scheduleUpdate(repoPath, channel)

	watcher.Watch(repoPath, channel)

	go func() {
		for {
			path := <-channel
			err = git.Sync(path)
			if err != nil {
				log.Printf("Syncing failed. Err: %v", err)
			}
		}
	}()

	log.Printf("Git notes is monitoring %s", repoPath)
}
