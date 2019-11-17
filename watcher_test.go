package main

import (
    "./test_helpers"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

type listener struct {
	paths []string
}

func setup() (Watcher, *listener, string) {
	var channel chan string = make(chan string)

	var watcher = GitWatcher {
		git: &GitCmd{},
		running: false,
		checkInterval: 10 * time.Millisecond,
		delayBeforeFiringEvent: 0,
		delayAfterFiringEvent: 1 * time.Second,
	}

	var path = test_helpers.SetupGitRepo("watcher")
	watcher.Watch(path, channel)

	var listener listener

	go func() {
		for {
			path = <- channel
			listener.paths = append(listener.paths, path)
		}
	}()

	return &watcher, &listener, path
}

func cleanup(watcher Watcher, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatalf("Unable to remove %s. Error: %v", path, err)
	}

	watcher.Stop()
}

func TestGitWatcher_CreateAndModify(t *testing.T) {
	var _, listener, path = setup()
	log.Println(path)

	//defer cleanup(watcher, path)

	assert.Equal(t, 0, len(listener.paths))

	test_helpers.WriteFile(t, path, "test.md", "Hello")
	time.Sleep(100 * time.Millisecond)
	log.Println(listener.paths)
	assert.Equal(t, 1, len(listener.paths))
	assert.Equal(t, path, listener.paths[0])

	test_helpers.WriteFile(t, path, "test.md", "Hello2")
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, 2, len(listener.paths))
	assert.Equal(t, path, listener.paths[0])
	assert.Equal(t, path, listener.paths[1])

	// No change
	test_helpers.WriteFile(t, path, "test.md", "Hello2")
	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, 2, len(listener.paths))
}

//func TestFsWatcher_Watch(t *testing.T) {
//	var listener, path = setup()
//	var currentCount = 0
//
//	// TODO: Add file
//
//	if listener.count <= currentCount {
//		t.Errorf("Count should be more than %d", currentCount)
//	}
//	currentCount = listener.count
//
//	// TODO: Modify file
//
//	if listener.count <= currentCount {
//		t.Errorf("Count should be more than %d", currentCount)
//	}
//	currentCount = listener.count
//
//	// TODO: Remove file
//
//	if listener.count <= currentCount {
//		t.Errorf("Count should be more than %d", currentCount)
//	}
//}
