package test_helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func SetupGitRepo(tag string, bare bool) string {
	path, err := ioutil.TempDir("", fmt.Sprintf("git_test_%s", tag))
	if err != nil {
		log.Fatalf("Unable to create a temp dir for the remote repo")
	}

	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}

	c := exec.Command("git", args...)
	c.Dir = path
	err = c.Run()
	if err != nil {
		log.Fatalf("Unable to init the remote repo")
	}

	return path
}

func SetupRemote(local string, remote string) {
	c := exec.Command("git", "remote", "add", "origin", remote)
	c.Dir = local
	err := c.Run()
	if err != nil {
		log.Fatalf("Unable to init the remote repo")
	}
}

func WriteFile(t *testing.T, repoPath string, filePath string, content string) {
	fullPath := fmt.Sprintf("%s/%s", repoPath, filePath)
	log.Printf("Write file: %v, content: %v", fullPath, content)
	assert.NoError(t, ioutil.WriteFile(fullPath, []byte(content), 0644))
}


func PerformCmd(t *testing.T, path string, cmd string, args... string) {
	log.Printf("Run cmd: %v", strings.Join(append([]string{cmd}, args...), " "))
	c := exec.Command(cmd, args...)
	c.Dir = path
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	assert.NoError(t, err)
}
