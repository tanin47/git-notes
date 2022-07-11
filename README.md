Git Notes
==========

[![CircleCI](https://circleci.com/gh/tanin47/git-notes.svg?style=svg)](https://circleci.com/gh/tanin47/git-notes)
[![codecov](https://codecov.io/gh/tanin47/git-notes/branch/master/graph/badge.svg)](https://codecov.io/gh/tanin47/git-notes)
[![Gitter](https://badges.gitter.im/tanin-git-notes/community.svg)](https://gitter.im/tanin-git-notes/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

*Your personal notes synced through Git*

Git Notes is in its alpha stage. I'd love to chat to users who want to use Git Notes. Please join [our Gitter channel](https://gitter.im/tanin-git-notes/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge).


------

Git Notes is a locally installed app that detects changes in a Git repo and push the changes immmediately to Github, Gitlab, or your own Git host. Here are the advantages:

* You can use your fav editor like Vim, Emacs, Sublime, or Atom.
* Your notes are more permanent. When was the last time you deleted a git repo? I don't remember mine either. Storing Github is how you're able to keep your several-year-old notes.
* Your notes are versioned by Git.
* Conflicts are handled intuitively for programmers. You see the git-style conflict text in your notes.

I hope Git Notes hits all the notes for you as it does for me. Enjoy!

  
Installation
-------------

0. Setup your personal note directory with Git. Make the master branch, commit, add `origin`, and `git push origin master -u`.
1. Install with `go install https://github.com/tanin47/git-notes@latest`
2. Make the config file that contains the paths that will be synced automatically by Git Notes. See the example: `git-notes.json.example`
	It can be found in `~/go/pkg/mod/github.com/tanin47/git-notes@{version with hash will be here}/`

The binary can be found as `git-notes` in the `~/go/bin`. (Or more accurate to say, in `$(go env GOBIN)` or even `$(go env GOPATH)/bin`)

You can run it by: `git-notes [your-config-file]`.
Probably you need to add next string to you `~/.profile`.
``` sh
export PATH=$(go env GOBIN):$PATH
```
And set GOBIN with `go env -w GOBIN=$(go env GOPATH)/bin`

To make Git Notes run at the startup and in the background, please follow the specific platform instruction below:

### Linux

Launch `git-notes --service`
This will add service file `~/.config/systemd/user/git-notes.service` and default config file in `~/.config/git-notes/git-notes.json`

Enable Git Notes to start at boot: `systemctl --user enable git-notes.service`

Run: `systemctl --user start git-notes.service`

Read logs: `journalctl --user -u git-notes.service --follow`


### Mac

Move `./service_conf/mac.git-notes.plist` to `~/Library/LaunchAgents/git-notes.plist`

Modify `~/Library/LaunchAgents/git-notes.plist` to use the binary that you built above with and your config file.

Run and start after booting:

1. `launchctl load ~/Library/LaunchAgents/git-notes.plist`
2. `launchctl start ~/Library/LaunchAgents/git-notes.plist`

If the plist file is changed, you will need to unload it first with: `launchctl unload ~/Library/LaunchAgents/git-notes.plist`.

Read logs: use Console.app. Search for `logger`.


### Windows

TBD

### Android and iOS

TBD: I want to build apps for this!

  
Architecture
-------------

Our main engine observes the current state of the git repo and make one action to transition to the next state.

Here are all the states:

* __dirty__: Unstaged change -> `git add .` -> __staged__
* __staged__: Staged change -> `git commit -m 'Updated'` -> __ahead__ or __out-of-sync__
* __ahead__: Ahead of the remote branch and can fast forward -> `git push` -> __synced__
* __out_of_sync__: The remote branch has unseen commits -> `git pull` -> __ahead__ (no conflict) or __dirty__ (there are conflicts)
* __synced__: The local branch matches the remote branch

This loop runs until no changes are observed. If the engine doesn't end on __synced__, something is wrong.

When the file change is detected, we invoke the engine again.

The file changes are detected by running `git status` every 10 seconds.

  
Develop
--------

* `go build` to build the binary
* `go run .` to run the application
* `go test` to run tests
* `gofmt -w .` to format all files
* `goimports -w .` to organize imports in all files
  

