package main

//import (
//	"testing"
//	"time"
//)
//
//func TestRun(t *testing.T) {
//	var watcher = MockWatcher{}
//	var git = MockGit{}
//
//	Run("some-path", &watcher, &git)
//
//	if watcher.repoPath != "some-path" {
//		t.Errorf("Path is not some-path")
//	}
//	if git.Count != 0 {
//		t.Errorf("Next invocation is not 0")
//	}
//
//	watcher.channel <- ""
//
//	// TODO: Figure out how to wait properly
//	time.Sleep(1 * time.Microsecond)
//	if git.Count != 1 {
//		t.Errorf("Next invocation is not 1")
//	}
//}
//
//type MockWatcher struct {
//	repoPath string
//	channel  chan string
//}
//
//func (m *MockWatcher) Watch(path string, channel chan string) {
//	m.repoPath = path
//	m.channel = channel
//}
//
//type MockGit struct {
//	Count int
//}
//
//func (m *MockGit) GetState(path string) State {
//	return Sync
//}
//
//func (m *MockGit) Next(path string) {
//	m.Count++
//}
