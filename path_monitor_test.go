package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGitRepoMonitor_StartMonitoring(t *testing.T) {
	gitRepoMonitor := GitRepoMonitor{
		scheduledUpdateInterval: time.Minute,
	}
	watcher := MockWatcher{}
	git := MockGit{}

	gitRepoMonitor.StartMonitoring("some-path", &watcher, &git)

	assert.Equal(t, "some-path", watcher.repoPath)
	assert.Equal(t, 1, git.Count)

	watcher.channel <- watcher.repoPath

	time.Sleep(1 * time.Second)
	assert.Equal(t, 2, git.Count)
}

func TestGitRepoMonitor_StartMonitoringAutomaticScheduleUpdate(t *testing.T) {
	gitRepoMonitor := GitRepoMonitor{
		scheduledUpdateInterval: 100 * time.Millisecond,
	}
	watcher := MockWatcher{}
	git := MockGit{}

	gitRepoMonitor.StartMonitoring("some-path", &watcher, &git)

	assert.Eventually(t, func() bool {
		return git.Count >= 2
	}, 1*time.Second, 10*time.Millisecond)
}

func TestGitRepoMonitor_ScheduleUpdate(t *testing.T) {
	gitRepoMonitor := GitRepoMonitor{
		scheduledUpdateInterval: 100 * time.Millisecond,
	}

	channel := make(chan string)
	var path string

	go func() {
		path = <-channel
	}()

	gitRepoMonitor.scheduleUpdate("some-path", channel)

	assert.Eventually(t, func() bool {
		return path == "some-path"
	}, 1*time.Second, 10*time.Millisecond)
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
