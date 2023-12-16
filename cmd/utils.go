package main

import (
	"os"
)

func checkBit(val uint64, bit uint) bool {
	if (uint(val) & bit) == 0 {
		return false
	}
	return true
}

func getPage(addr uint) uint {
	pagesize := uint(os.Getpagesize())
	return (addr & ^(pagesize - 1)) / pagesize
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
