package main

//import "testing"
//
//type listener struct {
//	count int
//}
//
//func setup() (listener, string) {
//	var channel chan string = make(chan string)
//
//	var watcher = NewFsWatcher()
//
//	var path = makeTempPath()
//	watcher.Watch(path, channel)
//
//	var listener listener
//
//	go func() {
//		for {
//			_ = <- channel
//			listener.count++
//		}
//	}()
//
//	return listener, path
//}
//
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
