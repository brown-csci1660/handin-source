package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("uid: ", os.Getuid())
	fmt.Println("euid:", os.Geteuid())
	fmt.Println("gid: ", os.Getgid())
	fmt.Println("egid:", os.Getegid())
}
