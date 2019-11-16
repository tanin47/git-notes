package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type repos struct {
	remotePath string
	remote git.Repository
	localPath  string
	local git.Repository
}

func setupRepos() repos {
	remotePath, err := ioutil.TempDir("", "git_test_remote")
	if err != nil {
		log.Fatalf("Unable to create a temp dir for the remote repo")
	}

	localPath, err := ioutil.TempDir("", "git_test_local")
	if err != nil {
		log.Fatalf("Unable to create a temp dir for the remote repo")
	}

	remote, err := git.PlainInit(remotePath, false)
	if err != nil {
		log.Fatalf("Unable to init the remote repo")
	}

	local, err := git.PlainInit(localPath, false)
	if err != nil {
		log.Fatalf("Unable to init the local repo")
	}

	_, err = local.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remotePath},
	})
	if err != nil {
		log.Fatalf("Unable to setup origin")
	}

	log.Printf("local: %s, remote: %s", localPath, remotePath)
	return repos{
		remotePath: remotePath,
		remote: *remote,
		localPath: localPath,
		local:  *local,
	}
}

func cleanupRepos(repos repos) {
	err := os.RemoveAll(repos.remotePath)
	if err != nil {
		log.Fatalf("Unable to remove %s. Error: %v", repos.remotePath, err)
	}

	err = os.RemoveAll(repos.localPath)
	if err != nil {
		log.Fatalf("Unable to remove %s. Error: %v", repos.localPath, err)
	}
}

func assertState(t *testing.T, path string, expectedState State) {
	gogit := GoGit{}
	state, err := gogit.GetState(path)
	assert.NoError(t, err)
	log.Printf("State: %v", state)
	assert.Equal(t, expectedState, state)
}

func performUpdate(t *testing.T, path string) {
	gogit := GoGit{}
	err := gogit.Update(path)
	assert.NoError(t, err)
}

func performSync(t *testing.T, path string) {
	gogit := GoGit{}
	err := gogit.Sync(path)
	assert.NoError(t, err)
}

func performCmd(t *testing.T, path string, cmd string, args... string) {
	log.Printf("Run cmd: %v", strings.Join(append([]string{cmd}, args...), " "))
	c := exec.Command(cmd, args...)
	c.Dir = path
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	assert.NoError(t, err)
}

func writeFile(t *testing.T, repoPath string, filePath string, content string) {
	fullPath := fmt.Sprintf("%s/%s", repoPath, filePath)
	log.Printf("Write file: %v, content: %v", fullPath, content)
	assert.NoError(t, ioutil.WriteFile(fullPath, []byte(content), 0644))
}

func TestGoGit_UpdateRename(t *testing.T) {
	repos := setupRepos()
	//defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test_name", "TestContent")

	assertState(t, repos.localPath, Dirty)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, Ahead)

	assert.NoError(t, os.Rename(fmt.Sprintf("%s/%s", repos.localPath, "test_name"), fmt.Sprintf("%s/%s", repos.localPath, "TEST_NAME")))

	assertState(t, repos.localPath, Dirty)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, OutOfSync)

	assertState(t, repos.localPath, OutOfSync)
	assertState(t, repos.localPath, OutOfSync)
}

func TestGoGit_UpdateModify(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")

	assertState(t, repos.localPath, Dirty)
	performSync(t, repos.localPath)
	assertState(t, repos.localPath, Sync)

	writeFile(t, repos.localPath, "test.md", "TestContent2")

	assertState(t, repos.localPath, Dirty)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, OutOfSync)
}

func TestGoGit_UpdateDirty(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")

	assertState(t, repos.localPath, Dirty)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, Ahead)
}

func TestGoGit_UpdateAhead(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test")

	assertState(t, repos.localPath, Ahead)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, Sync)
}

func TestGoGit_UpdateSync(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test")
	performCmd(t, repos.localPath, "git", "push", "origin", "master", "-u")

	assertState(t, repos.localPath, Sync)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, Sync)
}

func TestGoGit_UpdateOutOfSync(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test")
	performCmd(t, repos.localPath, "git", "push", "origin", "master", "-u")

	writeFile(t, repos.remotePath, "test.md", "UpdateFromRemote")
	performCmd(t, repos.remotePath, "git", "add", "--all")
	performCmd(t, repos.remotePath, "git", "commit", "-m", "Test")

	assertState(t, repos.localPath, OutOfSync)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, Sync)
}

func TestGoGit_UpdateFixConflict(t *testing.T) {
	repos := setupRepos()
	//defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test local")
	performCmd(t, repos.localPath, "git", "push", "origin", "master", "-u")

	writeFile(t, repos.remotePath, "test.md", "UpdateFromRemote")
	performCmd(t, repos.remotePath, "git", "add", "--all")
	performCmd(t, repos.remotePath, "git", "commit", "-m", "Test Remote")

	assertState(t, repos.localPath, OutOfSync)

	writeFile(t, repos.localPath, "test.md", "TestContent2")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test cause conflict")

	assertState(t, repos.localPath, OutOfSync)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, Dirty)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, Ahead)
	performUpdate(t, repos.localPath)
	assertState(t, repos.localPath, Sync)
}

func TestGoGit_SyncDirty(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")

	assertState(t, repos.localPath, Dirty)
	performSync(t, repos.localPath)
	assertState(t, repos.localPath, Sync)
}

func TestGoGit_SyncAhead(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test")

	assertState(t, repos.localPath, Ahead)
	performSync(t, repos.localPath)
	assertState(t, repos.localPath, Sync)
}

func TestGoGit_SyncSync(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test")
	performCmd(t, repos.localPath, "git", "push", "origin", "master", "-u")

	assertState(t, repos.localPath, Sync)
	performSync(t, repos.localPath)
	assertState(t, repos.localPath, Sync)
}

func TestGoGit_SyncOutOfSync(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test")
	performCmd(t, repos.localPath, "git", "push", "origin", "master", "-u")

	writeFile(t, repos.remotePath, "test.md", "UpdateFromRemote")
	performCmd(t, repos.remotePath, "git", "add", "--all")
	performCmd(t, repos.remotePath, "git", "commit", "-m", "Test")

	assertState(t, repos.localPath, OutOfSync)
	performSync(t, repos.localPath)
	assertState(t, repos.localPath, Sync)
}

func TestGoGit_SyncFixConflict(t *testing.T) {
	repos := setupRepos()
	defer cleanupRepos(repos)

	writeFile(t, repos.localPath, "test.md", "TestContent")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test local")
	performCmd(t, repos.localPath, "git", "push", "origin", "master", "-u")

	writeFile(t, repos.remotePath, "test.md", "UpdateFromRemote")
	performCmd(t, repos.remotePath, "git", "add", "--all")
	performCmd(t, repos.remotePath, "git", "commit", "-m", "Test Remote")

	assertState(t, repos.localPath, OutOfSync)

	writeFile(t, repos.localPath, "test.md", "TestContent2")
	performCmd(t, repos.localPath, "git", "add", "--all")
	performCmd(t, repos.localPath, "git", "commit", "-m", "Test cause conflict")

	assertState(t, repos.localPath, OutOfSync)
	performSync(t, repos.localPath)
	assertState(t, repos.localPath, Sync)
}

