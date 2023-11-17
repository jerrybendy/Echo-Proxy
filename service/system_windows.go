//go:build windows

package service

import (
	"os"
	"os/exec"
)

func isPrivileged() (bool, error) {
	return true, nil
}

func getHostsFilePath() string {
	r := os.Getenv("SYSTEMROOT")
	if r == "" {
		r = os.Getenv("windir")
	}
	return r + "\\system32\\drivers\\etc\\hosts"
}

func openFolder(path string) error {
	return exec.Command("explorer", path).Start()
}

func addHostsFileRecord() error {
	return addHostsRecordCommand()
}

func removeHostsFileRecord() error {
	return removeHostsRecordCommand()
}
