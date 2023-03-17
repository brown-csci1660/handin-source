// +build cs166 cs162

package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
)

var (
	handinDir = "/course/cs666/handin/"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %v [<assignment> [<tar-file>]]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	switch len(os.Args) {
	case 1:
		// there are no arguments; just show the usage and
		// list any available handins

		fmt.Printf("Usage: %v [<assignment> [<tar-file>]]\n\n", os.Args[0])
		files, err := ioutil.ReadDir(handinDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not list handins: %v\n", err)
			os.Exit(2)
		}
		var names []string
		for _, f := range files {
			if f.IsDir() {
				names = append(names, f.Name())
			}
		}
		if len(names) == 0 {
			fmt.Println("No available handins.")
			os.Exit(0)
		}
		sort.Strings(names)
		fmt.Println("Available handins:")
		for _, n := range names {
			fmt.Printf("\t%v\n", n)
		}
		os.Exit(0)
	case 2:
		// there is one argument; it's the name of the handin

		asgn := os.Args[1]
		hdir := filepath.Join(handinDir, asgn)

		// verify that the handin exists
		verifyHandin(hdir)

		// make sure we're in the current
		// user's home directory
		ok, err := sanitizePWD()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not verify pwd: %v\n", err)
			os.Exit(2)
		}
		if !ok {
			fmt.Fprintln(os.Stderr, "ERROR: Cannot hand in from outside your home directory!")
			os.Exit(1)
		}

		u := currentUser()
		userhdir := filepath.Join(hdir, u.Username)

		// create the user's handin directory
		// for this assignment if it doesn't
		// exist already
		makeUserHandinDir(userhdir)
		target := filepath.Join(userhdir, getRandomFilePrefix()+".tar")

		// tar the current directory and copy it to the
		// target file
		err = performTarHandin(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not hand in: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("successfully handed in to %v\n", target)
		fmt.Println("invoking autograder...")
		os.Exit(autograde(asgn, u, target))
	case 3:
		// there are two arguments; the first is the name of
		// the handin, and the second is the path to the tar

		asgn := os.Args[1]
		hdir := filepath.Join(handinDir, asgn)

		// verify that the handin exists
		verifyHandin(hdir)
		if len(os.Args[2]) < 4 || os.Args[2][len(os.Args[2])-4:] != ".tar" {
			fmt.Fprintln(os.Stderr, "path must end in .tar")
			os.Exit(1)
		}

		// make sure that the path is in the current
		// user's home directory
		tarpath := filepath.Clean(os.Args[2])
		tarpath, err := filepath.Abs(tarpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not make argument absolute path: %v\n", err)
			os.Exit(2)
		}
		ok, err := sanitizePath(tarpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not verify tar path: %v\n", err)
			os.Exit(2)
		}
		if !ok {
			fmt.Fprintln(os.Stderr, "ERROR: Cannot hand in file which is outside your home directory!")
			os.Exit(1)
		}

		u := currentUser()
		userhdir := filepath.Join(hdir, u.Username)

		// create the user's handin directory
		// for this assignment if it doesn't
		// exist already
		makeUserHandinDir(userhdir)

		// copy the named tar file into the handin directory
		target := filepath.Join(userhdir, getRandomFilePrefix()+".tar")
		t, err := os.Create(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create handin file: %v\n", err)
			os.Exit(2)
		}
		f, err := os.Open(tarpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open tar: %v\n", err)
			os.Exit(2)
		}
		_, err = io.Copy(t, f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not write handin file: %v\n", err)
			os.Exit(2)
		}
		err = t.Sync()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not sync handin file: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("successfully handed in to %v\n", target)
		fmt.Println("invoking autograder...")
		os.Exit(autograde(asgn, u, target))
	default:
		usage()
	}
}

func verifyHandin(hdir string) {
	fi, err := os.Stat(hdir)
	if (err != nil && os.IsNotExist(err)) || !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "no such handin: %v\n", hdir)
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "could not stat handin directory: %v\n", err)
		os.Exit(2)
	}
}

func currentUser() *user.User {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get current user: %v\n", err)
		os.Exit(2)
	}
	return u
}

func makeUserHandinDir(dir string) {
	err := os.MkdirAll(dir, 0770)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not make user handin directory: %v\n", err)
		os.Exit(2)
	}
}

func getRandomFilePrefix() string {
	var buf [16]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not generate random file prefix: %v\n", err)
		os.Exit(2)
	}
	return hex.EncodeToString(buf[:])
}
