// report is a setgid binary which allows
// students to see the current state of
// their grades in the course.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get current user: %v\n", err)
		os.Exit(2)
	}
	cmd := exec.Command("/course/cs666/tabin/modifydb", "--command", "view", "--student", u.Username)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running subcommand: %v\n", err)
		os.Exit(2)
	}
}
