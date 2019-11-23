package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanin47/git-notes/internal/test_helpers"
	"log"
	"os"
	"testing"
)


func assertState(t *testing.T, path string, expectedState State) {
	gogit := GitCmd{}
	state, err := gogit.GetState(path)
	assert.NoError(t, err)
	log.Printf("State: %v", state)
	assert.Equal(t, expectedState, state)
}

func performUpdate(t *testing.T, path string) {
	gogit := GitCmd{}
	err := gogit.Update(path)
	assert.NoError(t, err)
}

func performSync(t *testing.T, path string) {
	gogit := GitCmd{}
	err := gogit.Sync(path)
	assert.NoError(t, err)
}

func TestParseStatusBranch_NoRemote(t *testing.T) {
	state, err := ParseStatusBranch("## master")
	assert.NoError(t, err)
	assert.Equal(t, Ahead, state)
}

func TestParseStatusBranch_Sync(t *testing.T) {
	state, err := ParseStatusBranch("## master...origin/master")
	assert.NoError(t, err)
	assert.Equal(t, Sync, state)
}

func TestParseStatusBranch_Ahead(t *testing.T) {
	state, err := ParseStatusBranch("## master...origin/master [ahead 1]")
	assert.NoError(t, err)
	assert.Equal(t, Ahead, state)
}

func TestParseStatusBranch_OutOfSync(t *testing.T) {
	state, err := ParseStatusBranch("## master...origin/master [behind 99]")
	assert.NoError(t, err)
	assert.Equal(t, OutOfSync, state)
}

func TestParseStatusBranch_OutOfSync2(t *testing.T) {
	state, err := ParseStatusBranch("## master...origin/master [ahead 8, behind 99]")
	assert.NoError(t, err)
	assert.Equal(t, OutOfSync, state)
}

func TestGoGit_Rename(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test_name", "TestContent")

	assertState(t, repos.Local, Dirty)
	performSync(t, repos.Local)
	assertState(t, repos.Local, Sync)

	assert.NoError(t, os.Rename(fmt.Sprintf("%s/%s", repos.Local, "test_name"), fmt.Sprintf("%s/%s", repos.Local, "TEST_NAME")))

	assertState(t, repos.Local, Dirty)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Ahead)
}

func TestGoGit_Copy(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")

	assertState(t, repos.Local, Dirty)
	performSync(t, repos.Local)
	assertState(t, repos.Local, Sync)

	test_helpers.WriteFile(t, repos.Local, "copied.md", "TestContent")

	assertState(t, repos.Local, Dirty)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Ahead)
}

func TestGoGit_Modify(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")

	assertState(t, repos.Local, Dirty)
	performSync(t, repos.Local)
	assertState(t, repos.Local, Sync)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent2")

	assertState(t, repos.Local, Dirty)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Ahead)
}

func TestGoGit_Deletion(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")

	assertState(t, repos.Local, Dirty)
	performSync(t, repos.Local)
	assertState(t, repos.Local, Sync)

	assert.NoError(t, os.Remove(fmt.Sprintf("%s/%s", repos.Local, "test.md")))

	assertState(t, repos.Local, Dirty)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Ahead)
}

func TestGoGit_UpdateDirty(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")

	assertState(t, repos.Local, Dirty)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Ahead)
}

func TestGoGit_UpdateAhead(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test")

	assertState(t, repos.Local, Ahead)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Sync)
}

func TestGoGit_UpdateSync(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test")
	test_helpers.PerformCmd(t, repos.Local, "git", "push", "origin", "master", "-u")

	assertState(t, repos.Local, Sync)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Sync)
}

func TestGoGit_UpdateOutOfSync(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test")
	test_helpers.PerformCmd(t, repos.Local, "git", "push", "origin", "master", "-u")

	makeConflict(t, repos.Remote)

	assertState(t, repos.Local, OutOfSync)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Sync)
}

func TestGoGit_UpdateFixConflict(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test local")
	test_helpers.PerformCmd(t, repos.Local, "git", "push", "origin", "master", "-u")

	makeConflict(t, repos.Remote)
	assertState(t, repos.Local, OutOfSync)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent2")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test cause conflict")

	assertState(t, repos.Local, OutOfSync)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Dirty)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Ahead)
	performUpdate(t, repos.Local)
	assertState(t, repos.Local, Sync)
}

func TestGoGit_SyncDirty(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")

	assertState(t, repos.Local, Dirty)
	performSync(t, repos.Local)
	assertState(t, repos.Local, Sync)
}

func TestGoGit_SyncAhead(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test")

	assertState(t, repos.Local, Ahead)
	performSync(t, repos.Local)
	assertState(t, repos.Local, Sync)
}

func TestGoGit_SyncSync(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test")
	test_helpers.PerformCmd(t, repos.Local, "git", "push", "origin", "master", "-u")

	assertState(t, repos.Local, Sync)
	performSync(t, repos.Local)
	assertState(t, repos.Local, Sync)
}

func TestGoGit_SyncOutOfSync(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test")
	test_helpers.PerformCmd(t, repos.Local, "git", "push", "origin", "master", "-u")

	makeConflict(t, repos.Remote)

	assertState(t, repos.Local, OutOfSync)
	performSync(t, repos.Local)
	assertState(t, repos.Local, Sync)
}

func makeConflict(t *testing.T, remote string) {
	anotherLocal := test_helpers.SetupGitRepo("another_local", false)
	test_helpers.SetupRemote(anotherLocal, remote)
	test_helpers.PerformCmd(t, anotherLocal, "git", "fetch")
	test_helpers.PerformCmd(t, anotherLocal, "git", "checkout", "master")
	test_helpers.WriteFile(t, anotherLocal, "test.md", "Cause conflict")
	test_helpers.PerformCmd(t, anotherLocal, "git", "add", "--all")
	test_helpers.PerformCmd(t, anotherLocal, "git", "commit", "-m", "Test Remote")
	test_helpers.PerformCmd(t, anotherLocal, "git", "push")
}

func TestGoGit_SyncFixConflict(t *testing.T) {
	repos := test_helpers.SetupRepos()
	defer test_helpers.CleanupRepos(repos)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test local")
	test_helpers.PerformCmd(t, repos.Local, "git", "push", "origin", "master", "-u")

	makeConflict(t, repos.Remote)

	assertState(t, repos.Local, OutOfSync)

	test_helpers.WriteFile(t, repos.Local, "test.md", "TestContent2")
	test_helpers.PerformCmd(t, repos.Local, "git", "add", "--all")
	test_helpers.PerformCmd(t, repos.Local, "git", "commit", "-m", "Test cause conflict")

	assertState(t, repos.Local, OutOfSync)
	performSync(t, repos.Local)
	assertState(t, repos.Local, Sync)
}

