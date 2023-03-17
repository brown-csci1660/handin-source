// +build !cs162,cs166

package main

import (
	"os"
	"os/exec"
	"os/user"
)

func performTarHandin(target string) error {
	cmd := exec.Command("/bin/tar", "-hcvf", target, ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// sanitizePath ensures that the given path
// is in the current user's home directory
func sanitizePath(path string) (ok bool, err error) {
	u, err := user.Current()
	if err != nil {
		return false, err
	}

	return stringHasPrefix(path, u.HomeDir), nil
}

// sanitizePWD ensures that the current working
// directory is in the current user's home directory
func sanitizePWD() (ok bool, err error) { return sanitizePath(os.Getenv("PWD")) }

func stringHasPrefix(s, prefix string) bool {
	l := len(prefix)
	if len(s) < l {
		return false
	}
	return s[:l] == prefix
}
