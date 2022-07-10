package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	Error     State = "error"
	Dirty     State = "dirty"
	Ahead     State = "ahead"
	OutOfSync State = "out-of-sync"
	Sync      State = "sync"
)

type State string

type Git interface {
	IsDirty(path string) (bool, error)
	GetState(path string) (State, error)
	Sync(path string) error
	Update(path string) error
}

type GitCmd struct{}

func (g *GitCmd) Sync(path string) error {
	state, err := g.GetState(path)
	log.Printf("Starting state: %s", state)
	if err != nil {
		return fmt.Errorf("performing GetState() failed. Err: %v", err)
	}

	for {
		if state == Sync {
			return nil
		}

		err = g.Update(path)
		if err != nil {
			return fmt.Errorf("performing Update() failed. Err: %v", err)
		}
		nextState, err := g.GetState(path)
		if err != nil {
			return fmt.Errorf("performing GetState() failed. Err: %v", err)
		}
		log.Printf("Next state: %s", nextState)

		if state == nextState {
			return fmt.Errorf("state doesn't change. Something is wrong")
		}

		state = nextState
	}
}

func runCmd(path string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = path

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (g *GitCmd) IsDirty(path string) (bool, error) {
	out, err := runCmd(path, "git", "status", "--porcelain")
	if err != nil {
		return false, fmt.Errorf("unable to get status. Error: %v", err)
	}

	dirty := strings.TrimSpace(string(out)) != ""
	return dirty, nil
}

func (g *GitCmd) GetState(path string) (State, error) {
	log.Printf("Computing the state of %s", path)

	dirty, err := g.IsDirty(path)
	if err != nil {
		return Error, fmt.Errorf("unable to get status. Error: %v", err)
	}
	if dirty {
		return Dirty, nil
	} else {
		state, err := GetStateAgainstRemote(path)
		if err != nil {
			return Error, err
		}
		return state, nil
	}
}

func ParseStatusBranch(status string) (State, error) {
	// 5 variants of status branch
	// ## master
	// ## master...origin/master
	// ## master...origin/master [ahead 1]
	// ## master...origin/master [behind 1]
	// ## master...origin/master [ahead 1, behind 1]

	reg := regexp.MustCompile("## master(\\.\\.\\.origin\\/master *(\\[(ahead|behind) *[0-9]+ *(, *behind *[0-9]+)? *])?)?")
	matches := reg.FindAllStringSubmatch(status, -1)

	if len(matches) == 0 {
		return Error, fmt.Errorf("unable to parse status: %v", status)
	}

	groups := matches[0]
	if groups[0] == "" {
		return Error, fmt.Errorf("unable to parse status: %v", status)
	}

	// ## master
	if groups[1] == "" {
		return Ahead, nil
	}

	// ## master...origin/master
	if groups[2] == "" {
		return Sync, nil
	}

	// ## master...origin/master
	if groups[3] == "behind" {
		return OutOfSync, nil
	}

	// ## master...origin/master
	if groups[3] == "ahead" {
		if groups[4] == "" {
			return Ahead, nil
		} else {
			return OutOfSync, nil
		}
	}

	return Error, fmt.Errorf("unable to parse status: %v", status)
}

func GetStateAgainstRemote(path string) (State, error) {
	_, err := runCmd(path, "git", "fetch")
	if err != nil {
		return Error, fmt.Errorf("unable to fetch. Error: %v", err)
	}

	status, err := runCmd(path, "git", "status", "--branch", "--porcelain")
	if err != nil {
		return Error, fmt.Errorf("unable to fetch. Error: %v", err)
	}

	return ParseStatusBranch(status)
}

func (g *GitCmd) Update(path string) error {
	state, err := g.GetState(path)
	if err != nil {
		return err
	}

	switch state {
	case Error:
	case Dirty:
		err = AddAndCommit(path)
	case Ahead:
		err = Push(path)
	case OutOfSync:
		err = Merge(path)
	case Sync:
	}

	return err
}

func AddAndCommit(path string) error {
	err := Add(path)
	if err != nil {
		return err
	}
	return Commit(path)
}

func Merge(path string) error {
	cmd := exec.Command("git", "merge", "origin/master", "--allow-unrelated-histories", "--no-commit")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // Merge fails if there's conflict. So, we ignore the failure.
	return nil
}

func Push(path string) error {
	cmd := exec.Command("git", "push", "origin", "master", "-u")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Add(path string) error {
	cmd := exec.Command("git", "add", "--all")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Commit(path string) error {
	cmd := exec.Command("git", "-c", "user.name='Git notes'", "-c", "user.email='git-notes@noemail.com'", "commit", "-m", fmt.Sprintf("Commited at %v", time.Now()))
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NewGoGit() GitCmd {
	return GitCmd{}
}
