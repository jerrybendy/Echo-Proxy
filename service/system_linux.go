//go:build linux

package service

import (
	"os/exec"
	"os/user"
)

func isPrivileged() (bool, error) {
	u, err := user.Current()
	if err != nil {
		return false, err
	}
	return u.Uid == "0", nil
}

func getHostsFilePath() string {
	return "/etc/hosts"
}

func openFolder(path string) error {
	return exec.Command("open", path).Start()
}

func addHostsFileRecord() error {
	return addHostsRecordCommand()
}

func removeHostsFileRecord() error {
	return removeHostsRecordCommand()
}
