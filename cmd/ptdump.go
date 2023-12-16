package main

import (
	"flag"
	"strconv"
)

var granular_flag bool = false

func main() {
	init_flags()
	flag.Parse()
	args := flag.Args()

	pid, err := strconv.Atoi(args[0])
	check(err)

	pt := NewPageTable(pid)
	pt.Dump()
	/*
		fmt.Printf("%v%30v%13v%15v%15v%15v%11v\n", "virt addr", "physical addr", "size", "perms", "present", "swapped", "path")
		for _, pte := range pt.Entries() {
			m, pageinfo := pte.Unpack()
			fmt.Printf("%0#x%#20x%20d%12s%13t%16t\t\t%s\n",
				m.StartAddr, pageinfo.Addr, (m.EndAddr - m.StartAddr),
				m.Perms, pageinfo.Present, pageinfo.IsSwapped, m.Path)
		}
	*/
}

func init_flags() {
	flag.BoolVar(&granular_flag, "g", false, "Granular output: show all individual pages, instead of them being coaleced")
}
