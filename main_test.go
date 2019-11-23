package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun(t *testing.T) {
	var git = MockGit{}
	var watcher = MockWatcher{}
	var configReader = MockConfigReader{}
	var monitor = MockMonitor{}

	Run(&git, &watcher, &configReader, &monitor)

	assert.Equal(t, []string{"some-path", "some-path-2"}, monitor.startMonitorPaths)
}

type MockConfigReader struct {}

func (m *MockConfigReader) Read(path string) (*Config, error) {
	var config = &Config{
		Repos: []string{"some-path", "some-path-2"},
	}
	return config, nil
}

type MockMonitor struct {
	startMonitorPaths []string
}

func (m *MockMonitor) StartMonitoring(repoPath string, watcher Watcher, git Git) {
	m.startMonitorPaths = append(m.startMonitorPaths, repoPath)
}

func (m *MockMonitor) scheduleUpdate(repoPath string, channel chan string) {
}



