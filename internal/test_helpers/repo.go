package test_helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Repos struct {
	Remote string
	Local  string
}

func SetupRepos() Repos {
	remote := SetupGitRepo("Remote", true)
	local := SetupGitRepo("Local", false)

	SetupRemote(local, remote)

	log.Printf("Local: %s, Remote: %s", local, remote)
	return Repos{
		Remote: remote,
		Local:  local,
	}
}

func CleanupRepos(repos Repos) {
	err := os.RemoveAll(repos.Remote)
	if err != nil {
		log.Fatalf("Unable to remove %s. Error: %v", repos.Remote, err)
	}

	err = os.RemoveAll(repos.Local)
	if err != nil {
		log.Fatalf("Unable to remove %s. Error: %v", repos.Local, err)
	}
}

func SetupGitRepo(tag string, bare bool) string {
	path, err := ioutil.TempDir("", fmt.Sprintf("git_test_%s", tag))
	if err != nil {
		log.Fatalf("Unable to create a temp dir for the Remote repo")
	}

	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}

	c := exec.Command("git", args...)
	c.Dir = path
	err = c.Run()
	if err != nil {
		log.Fatalf("Unable to init the repo. Path: %v, Error: %v", path, err)
	}

	return path
}

func SetupRemote(local string, remote string) {
	c := exec.Command("git", "remote", "add", "origin", remote)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Dir = local
	err := c.Run()
	if err != nil {
		log.Fatalf("Unable to add origin. Error: %v", err)
	}
}

func WriteFile(t *testing.T, repoPath string, filePath string, content string) {
	fullPath := fmt.Sprintf("%s/%s", repoPath, filePath)
	log.Printf("Write file: %v, content: %v", fullPath, content)
	assert.NoError(t, ioutil.WriteFile(fullPath, []byte(content), 0o644))
}

func PerformCmd(t *testing.T, path string, cmd string, args ...string) {
	log.Printf("Run cmd: %v", strings.Join(append([]string{cmd}, args...), " "))
	c := exec.Command(cmd, args...)
	c.Dir = path
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	assert.NoError(t, err)
}
