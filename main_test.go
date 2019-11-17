package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	var watcher = MockWatcher{}
	var git = MockGit{}

	Run("some-path", &watcher, &git)

	assert.Equal(t, "some-path", watcher.repoPath)
	assert.Equal(t, 1, git.Count)

	watcher.channel <- watcher.repoPath

	time.Sleep(1 * time.Second)
	assert.Equal(t, 2, git.Count)
}

type MockWatcher struct {
	repoPath string
	channel  chan string
}

func (m *MockWatcher) Watch(path string, channel chan string) {
	m.repoPath = path
	m.channel = channel
}

type MockGit struct {
	Count int
}

func (m *MockGit) IsDirty(path string) (bool, error) {
	return false, nil
}

func (m *MockGit) Sync(path string) error {
	m.Count++
	return nil
}

func (m *MockGit) Update(path string) error {
	return nil
}

func (m *MockGit) GetState(path string) (State, error) {
	return Sync, nil
}
