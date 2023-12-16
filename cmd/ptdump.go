package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s pid\n", os.Args[0])
		os.Exit(1)
	}

	pid, err := strconv.Atoi(os.Args[1])
	check(err)

	maps := getMaps(pid)
	fmt.Printf("%v%30v%13v%15v%15v%15v%11v\n", "virt addr", "physical addr", "size", "perms", "present", "swapped", "path")
	for _, m := range maps {
		pageinfo := getPageInfo(pid, m.StartAddr)
		fmt.Printf("%0#x%#20x%20d%12s%13t%16t\t\t%s\n",
			m.StartAddr, pageinfo.Addr, (m.EndAddr - m.StartAddr),
			m.Perms, pageinfo.Present, pageinfo.IsSwapped, m.Path)
	}
}
