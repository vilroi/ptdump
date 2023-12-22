package pagetable

func checkBit(val uint64, bit uint) bool {
	if (uint(val) & bit) == 0 {
		return false
	}
	return true
}

func getPage(addr uint) uint {
	return (addr & ^(pagesize - 1)) / pagesize
}

func pageAlign(addr uint) uint {
	return (addr & ^(pagesize - 1))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
