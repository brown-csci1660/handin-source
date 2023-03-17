package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"

	"../lib"
)

// autograde extracts the handin and runs the autograder,
// which both grades the handin and inserts the grades
// into the database.
func autograde(asgn string, u *user.User, file string) (code int) {
	// create a temp directory to extract into
	tdir, err := ioutil.TempDir("", "cs666_autograde")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create directory for autograding: %v\n", err)
		return 2
	}
	defer os.RemoveAll(tdir)

	// extract the handin into the temp directory
	err = lib.ExtractTar(file, tdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not extract handin: %v\n", err)
		return 2
	}

	err = os.Chdir(tdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not cd: %v\n", err)
		return 2
	}

	// err = syscall.Setgid(os.Getegid())
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "could not raise group privileges: %v\n", err)
	// 	return 2
	// }

	// run the autograder
	cmd := exec.Command("/course/cs666/tabin/autograde.sh", asgn, u.Username)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not run autograder: %v\n", err)
		return 2
	}
	return 0
}
