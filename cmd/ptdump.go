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
}

func init_flags() {
	flag.BoolVar(&granular_flag, "g", false, "Granular output: show all individual pages, instead of them being coaleced")
}
