package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	Error State = "error"
	Dirty  State = "dirty"
	Ahead     State = "ahead"
	OutOfSync State = "out-of-sync"
	Sync      State = "sync"
)

type State string

type Git interface {
	GetState(path string) (State, error)
	Sync(path string) error
	Update(path string) error
}

type GoGit struct {
}

func GetRepoAndWorktree(path string) (*git.Repository, *git.Worktree, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open %s. Error: %v", path, err)
	}

	current, err := repo.Worktree()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get the current working tree. Error: %v", err)
	}

	return repo, current, nil
}

func (g *GoGit) Sync(path string) error {
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

func (g *GoGit) GetState(path string) (State, error) {
	log.Printf("Computing the state of %s", path)

	repo, current, err := GetRepoAndWorktree(path)
	if err != nil {
		return Error, fmt.Errorf("unable to get the repo and worktree. Error: %v", err)
	}

	status, err := current.Status()
	if err != nil {
		return Error, fmt.Errorf("unable to get the status. Error: %v", err)
	}

	log.Println(status)

	if len(status) > 0 {
		return Dirty, nil
	} else {
		state, err := GetStateAgainstRemote(*repo, *current)
		if err != nil {
			return Error, err
		}
		return state, nil
	}
}

func GetStateAgainstRemote(repo git.Repository, current git.Worktree) (State, error) {
	err := repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"refs/heads/master:refs/remotes/origin/master"},
	})

	if err != nil {
		if err == transport.ErrEmptyRemoteRepository || strings.Contains(err.Error(), "couldn't find remote ref") {
			return Ahead, nil
		} else if err == git.NoErrAlreadyUpToDate || strings.Contains(err.Error(), "reference not found"){
			// Do nothing
		} else {
			return Error, fmt.Errorf("unable to fetch. Error: %v", err)
		}
	}

	remoteRef, err := repo.Reference("refs/remotes/origin/master", true)
	if err != nil {
		return Error, fmt.Errorf("unable to get the remote reference refs/remotes/origin/master. Error: %v", err)
	}
	remoteCommit, err := repo.CommitObject(remoteRef.Hash())
	if err != nil {
		return Error, fmt.Errorf("unable to get the commit %s, which corresponds to remotes/origin/master. Error: %v", remoteRef.Hash().String(), err)
	}

	localCommitIter, err := repo.Log(&git.LogOptions{
		From:     plumbing.Hash{},
		Order:    0,
		FileName: nil,
		All:      false,
	})
	if err != nil {
		return Error, fmt.Errorf("unable to get the list of local commits. Error: %v", err)
	}
	localCommit, err := localCommitIter.Next()
	if err != nil {
		return Error, fmt.Errorf("unable to get the latest local commit. Error: %v", err)
	}
	if remoteCommit.Hash == localCommit.Hash {
		return Sync, nil
	} else {
		isAncestor, err := remoteCommit.IsAncestor(localCommit)
		if err != nil {
			return Error, fmt.Errorf("unable to check if the remote commit is the ancestor of the local commit. Error: %v", err)
		}

		if isAncestor {
			return Ahead, nil
		} else {
			return OutOfSync, nil
		}
	}
}

func (g *GoGit) Update(path string) error {
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

func NewGoGit() GoGit {
	return GoGit{}
}
