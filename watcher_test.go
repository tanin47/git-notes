package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/tanin47/git-notes/internal/test_helpers"
	"log"
	"os"
	"testing"
	"time"
)

type listener struct {
	paths []string
}

func setup() (*GitWatcher, *listener, string, chan string) {
	var channel chan string = make(chan string)

	var watcher = GitWatcher {
		git: &GitCmd{},
		running: false,
		checkInterval: 10 * time.Millisecond,
		delayBeforeFiringEvent: 0,
		delayAfterFiringEvent: 1 * time.Second,
	}

	var path = test_helpers.SetupGitRepo("watcher", false)

	var listener listener

	go func() {
		for {
			path = <- channel
			listener.paths = append(listener.paths, path)
		}
	}()

	return &watcher, &listener, path, channel
}

func cleanup(watcher *GitWatcher, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatalf("Unable to remove %s. Error: %v", path, err)
	}

	watcher.Stop()
}

func commit(t *testing.T, path string) {
	test_helpers.PerformCmd(t, path, "git", "add", "--all")
	test_helpers.PerformCmd(t, path, "git", "commit", "-m", "Test")
}

func TestGitWatcher_Watch(t *testing.T) {
	var watcher, listener, path, channel = setup()
	defer cleanup(watcher, path)

	watcher.Watch(path, channel)

	assert.Equal(t, 0, len(listener.paths))

	test_helpers.WriteFile(t, path, "test.md", "Watch")
	time.Sleep(1 * time.Second)
	assert.Greater(t, len(listener.paths), 0)
	assert.Equal(t, path, listener.paths[0])
}

func TestGitWatcher_CreateAndModify(t *testing.T) {
	var watcher, listener, path, channel = setup()
	defer cleanup(watcher, path)

	watcher.Check(path, channel)
	assert.Equal(t, 0, len(listener.paths))

	test_helpers.WriteFile(t, path, "test.md", "Hello")
	watcher.Check(path, channel)
	assert.Equal(t, 1, len(listener.paths))
	assert.Equal(t, path, listener.paths[0])

	commit(t, path)

	watcher.Check(path, channel)
	assert.Equal(t, 1, len(listener.paths))
	assert.Equal(t, path, listener.paths[0])

	test_helpers.WriteFile(t, path, "test.md", "Hello2")
	watcher.Check(path, channel)
	assert.Equal(t, 2, len(listener.paths))
	assert.Equal(t, path, listener.paths[0])
	assert.Equal(t, path, listener.paths[1])

	commit(t, path)

	// No change
	test_helpers.WriteFile(t, path, "test.md", "Hello2")
	watcher.Check(path, channel)
	assert.Equal(t, 2, len(listener.paths))
}
