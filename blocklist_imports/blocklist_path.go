// +build cs166,!cs162

package main

func blocklistPath(path string) bool {
	return blocklistMap[path]
}

var blocklist = []string{
	"flag",
	"fmt",
	"io/ioutil",
	"net",
	"net/http",
	"net/rpc",
	"net/smtp",
	"os",
	"os/exec",
	"syscall",
	"unsafe",
}

func init() {
	blocklistMap = make(map[string]bool)
	for _, b := range blocklist {
		blocklistMap[b] = true
	}
}

var blocklistMap map[string]bool
