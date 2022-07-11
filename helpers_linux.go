//go:build linux

package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"

	_ "embed"
)

//go:embed service_conf/linux.git-notes.service
var serviceData []byte

// makeService installs a service file in $XDG_DATA_HOME/systemd/user or $HOME/.local/share/systemd/user
// More information can be found on https://www.freedesktop.org/software/systemd/man/systemd.unit.html
func makeService() {
	var err error
	home, err := os.UserHomeDir()

	var pathToService string
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		log.Println("Creating $XDG_CONFIG_HOME/systemd/user/git-notes.service")
		pathToService = path.Join(os.Getenv("XDG_CONFIG_HOME"), "systemd/user")
	} else {
		log.Println("Creating ~/.config/systemd/user/git-notes.service")
		pathToService = path.Join(home, ".config/systemd/user")
	}

	if err = os.MkdirAll(pathToService, os.ModePerm); err != nil {
		log.Fatalln("Failed making directory for service files", err)
		return
	}

	pathToService = path.Join(pathToService, "git-notes.service")
	if err = writeService(pathToService, serviceData); err != nil {
		if os.IsExist(err) {
			log.Println("Config already exist, delete it manually if you want to create new one")
		} else {
			log.Fatal("Error writing service file", err)
		}
	}

	prc := exec.Command("systemctl", "--user", "daemon-reload")
	if _, err = prc.CombinedOutput(); err != nil {
		log.Println(prc.String())
		log.Println(err)
	}
}

func makeConfig() {
	var err error
	var home string
	var configData []byte

	home, err = os.UserHomeDir()

	configData, err = json.MarshalIndent(
		Config{
			make([]string, 0),
		},
		"", "\t",
	)

	var pathToConfig string
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		log.Println("Creating $XDG_CONFIG_HOME/git-notes/git-notes.json")
		pathToConfig = path.Join(os.Getenv("XDG_CONFIG_HOME"), "git-notes")
	} else {
		log.Println("Creating ~/.config/git-notes/git-notes.json")
		pathToConfig = path.Join(home, ".config/git-notes")
	}

	if err = os.MkdirAll(pathToConfig, os.ModePerm); err != nil {
		log.Fatalln("Failed making directory for config files", err)
		return
	}

	pathToConfig = path.Join(pathToConfig, "git-notes.json")
	if err = writeService(pathToConfig, configData); err != nil {
		if os.IsExist(err) {
			log.Println("Config already exist, delete it manually if you want to create new one")
		} else {
			log.Fatal("Error writing service file", err)
		}
	}
}

func writeService(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}

	if _, err = f.Write(data); err != nil {
		return err
	}

	if err = f.Sync(); err != nil {
		return err
	}

	return f.Close()
}
