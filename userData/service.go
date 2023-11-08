package userData

import "localProxy/hostsFile"

type Service struct {
}

func (s *Service) StartServer() bool {
	isPrivileged, err := hostsFile.IsPrivileged()
	if err != nil {
		emitErrorToFrontend(err.Error())
		return false
	}

	if !isPrivileged {
		emitErrorToFrontend("You are not running with root/administrator permissions, `apply to /etc/hosts` is skipped!")
	}

	return true
}
