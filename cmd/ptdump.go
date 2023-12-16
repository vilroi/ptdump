package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/vilroi/ptdump/internal/pagetable"
)

func main() {
	init_flags()

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	pid, err := strconv.Atoi(args[0])
	check(err)

	pt := pagetable.NewPageTable(pid)
	pt.Dump()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func init_flags() {
	flag.BoolVar(&pagetable.GranularFlag, "g", false, "Granular output: show all individual pages, instead of them being coaleced")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [opts] pid\n", os.Args[0])
		flag.PrintDefaults()
	}
}
