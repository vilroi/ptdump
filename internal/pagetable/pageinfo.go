package pagetable

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	PRESENT_BIT    uint = (1 << 63)
	SWAP_BIT       uint = (1 << 62)
	MAP_BIT        uint = (1 << 61)
	WRITE_PROT_BIT uint = (1 << 57)
	EXCL_MAP_BIT   uint = (1 << 56)
	SOFT_DIRTY_BIT uint = (1 << 55)
)

const (
	PADDR_MASK uint64 = (1 << 54) - 1
)

type PageInfo struct {
	present      bool
	isSwapped    bool
	isFileOrAnon bool
	writeProt    bool
	mapExcl      bool
	softDirty    bool
	addr         uint
}

// TODO: change from hardcoded value to something arch dependent
// (8 bytes assumes 64 bit?)
func getPageInfo(pid int, virt_addr uint) PageInfo {
	path := fmt.Sprintf("/proc/%d/pagemap", pid)
	f, err := os.Open(path)
	check(err)

	page := getPage(virt_addr)
	offset := int64(page * 8)
	_, err = f.Seek(offset, os.SEEK_SET)
	check(err)

	buf := make([]byte, 8)
	_, err = f.Read(buf)
	check(err)

	pte := binary.LittleEndian.Uint64(buf)

	return newPageInfo(pte)
}

func newPageInfo(pte uint64) PageInfo {
	var pageinfo PageInfo

	pageinfo.present = checkBit(pte, PRESENT_BIT)
	pageinfo.isSwapped = checkBit(pte, SWAP_BIT)
	pageinfo.isFileOrAnon = checkBit(pte, MAP_BIT)
	pageinfo.writeProt = checkBit(pte, WRITE_PROT_BIT)
	pageinfo.mapExcl = checkBit(pte, EXCL_MAP_BIT)
	pageinfo.softDirty = checkBit(pte, SOFT_DIRTY_BIT)
	pageinfo.addr = uint(int(pte&PADDR_MASK) * os.Getpagesize())

	return pageinfo
}
