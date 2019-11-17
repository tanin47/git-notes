package test_helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os/exec"
	"testing"
)

func SetupGitRepo(tag string) string {
	path, err := ioutil.TempDir("", fmt.Sprintf("git_test_%s", tag))
	if err != nil {
		log.Fatalf("Unable to create a temp dir for the remote repo")
	}

	c := exec.Command("git", "init")
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

