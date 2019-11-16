package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Repos []string `json:"repos"`
}

func main() {
	log.Println("Git notes is starting...")

	var git = NewGoGit()
	var watcher = GitWatcher{
		git:     &git,
		running: false,
	}

	if len(os.Args) < 2 {
		log.Fatal("Please pass the config file path as the first argument.")
	}
	configPath := os.Args[1]

	config, err := ReadConfig(configPath)

	if err != nil {
		log.Fatalf("Unable to read the config file. Err: %v", err)
	}

	fmt.Println(config)
	for _, repoPath := range config.Repos {
		Run(repoPath, &watcher, &git)
	}
	select {}
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {  return nil, err }

	decoder := json.NewDecoder(file)

	var config Config
	err = decoder.Decode(&config)
	if err != nil {  return nil, err }

	return &config, nil
}

func ScheduleUpdate(channel chan string) {
	time.AfterFunc(1*time.Hour, func() {
		channel <- ""
		ScheduleUpdate(channel)
	})
}

func Run(repoPath string, watcher Watcher, git Git) {
	var channel = make(chan string)
	err := git.Sync(repoPath)
	if err != nil {
		log.Printf("Syncing failed. Err: %v", err)
	}
	ScheduleUpdate(channel)

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
