package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Running = true

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--service" {
		log.Println("Git Notes start creating service...")
		makeService()
		makeConfig()
		os.Exit(0)
	}

	log.Println("Git Notes is starting...")

	git := NewGoGit()
	watcher := GitWatcher{
		git:                    &git,
		running:                false,
		checkInterval:          10 * time.Second,
		delayBeforeFiringEvent: 2 * time.Second,
		delayAfterFiringEvent:  5 * time.Second,
	}
	configReader := JsonConfigReader{}
	gitRepoMonitor := GitRepoMonitor{
		scheduledUpdateInterval: 5 * time.Minute,
	}

	Run(&git, &watcher, &configReader, &gitRepoMonitor)

	for Running {
		time.Sleep(1 * time.Second)
	}
}

func Run(git Git, watcher Watcher, configReader ConfigReader, monitor PathMonitor) {
	if len(os.Args) < 2 {
		log.Fatal("Please pass the config file path as the first argument.")
	}
	configPath := os.Args[1]
	config, err := configReader.Read(configPath)
	if err != nil {
		log.Fatalf("Unable to read the config file. Err: %v", err)
	}

	fmt.Println(config)
	for _, repoPath := range config.Repos {
		monitor.StartMonitoring(repoPath, watcher, git)
	}
}
