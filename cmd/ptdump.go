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
	fmt.Println("virt Addr\t\tphysical addr\t\tsize\t\tperms\t\tpresent\t\tswapped")
	for _, m := range maps {
		pageinfo := getPageInfo(pid, m.StartAddr)
		fmt.Printf("0x%x\t\t0x%x\t\t%d\t\t%s\t\t%t\t\t%t\t\t%s\n",
			m.StartAddr, pageinfo.Addr, (m.EndAddr - m.StartAddr),
			m.Perms, pageinfo.Present, pageinfo.IsSwapped, m.Path)
	}
}
