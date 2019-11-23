package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanin47/git-notes/internal/test_helpers"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestMainFunc(t *testing.T) {
	Running = true

	var git = NewGoGit()

	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	configDir, err := ioutil.TempDir("", "git-notes-config-dir")
	assert.NoError(t, err)

	test_helpers.WriteFile(t, configDir, "git-notes.json", fmt.Sprintf(`{ "repos": [ "%s" ] }`, repos.Local))

	oldArgs := os.Args
	os.Args = []string{"app", fmt.Sprintf("%s/%s", configDir, "git-notes.json")}
	defer func() { os.Args = oldArgs }()

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "First commit")

	state, err := git.GetState(repos.Local)
	assert.NoError(t, err)
	assert.Equal(t, Ahead, state)

	go main()

	assert.Eventually(t, func() bool {
		state, err := git.GetState(repos.Local)
		assert.NoError(t, err)
		return state == Sync
	}, 15 * time.Second, 1 * time.Second)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent2")

	state, err = git.GetState(repos.Local)
	assert.NoError(t, err)
	assert.Equal(t, Dirty, state)

	assert.Eventually(t, func() bool {
		state, err := git.GetState(repos.Local)
		assert.NoError(t, err)
		return state == Sync
	}, 15 * time.Second, 1 * time.Second)

	Running = false
}

func TestRun(t *testing.T) {
	var git = MockGit{}
	var watcher = MockWatcher{}
	var configReader = MockConfigReader{}
	var monitor = MockMonitor{}

	oldArgs := os.Args
	os.Args = []string{"app", "some-git-notes.json"}
	defer func() { os.Args = oldArgs }()

	Run(&git, &watcher, &configReader, &monitor)

	assert.Equal(t, "some-git-notes.json", configReader.readPath)
	assert.Equal(t, []string{"some-path", "some-path-2"}, monitor.startMonitorPaths)
}

type MockConfigReader struct {
	readPath string
}

func (m *MockConfigReader) Read(path string) (*Config, error) {
	m.readPath = path
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
