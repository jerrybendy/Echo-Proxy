//go:build darwin

package service

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
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
	if p, _ := isPrivileged(); p {
		return addHostsRecordCommand()
	}
	// Only runs when ApplyHosts is enabled
	isEnabled := false
	for _, h := range config.Hosts {
		if h.ApplyHosts {
			isEnabled = true
			break
		}
	}
	if isEnabled {
		cmd := strings.ReplaceAll(os.Args[0], " ", "\\\\ ") + " addHostsRecord"
		return exec.Command("osascript", "-e", fmt.Sprintf("do shell script \"%s\" with administrator privileges", cmd)).Start()
	}
	return nil
}

func removeHostsFileRecord() error {
	if p, _ := isPrivileged(); p {
		return removeHostsRecordCommand()
	}
	// Only runs when ApplyHosts is enabled
	isEnabled := false
	for _, h := range config.Hosts {
		if h.ApplyHosts {
			isEnabled = true
			break
		}
	}
	if isEnabled {
		cmd := strings.ReplaceAll(os.Args[0], " ", "\\\\ ") + " removeHostsRecord"
		return exec.Command("osascript", "-e", fmt.Sprintf("do shell script \"%s\" with administrator privileges", cmd)).Start()
	}
	return nil
}
