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
	Present      bool
	IsSwapped    bool
	IsFileOrAnon bool
	WriteProt    bool
	MapExcl      bool
	SoftDirty    bool
	Addr         uint
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
	//fmt.Printf("0b%b\n", pte)

	return newPageInfo(pte)
}

func newPageInfo(pte uint64) PageInfo {
	var pageinfo PageInfo

	pageinfo.Present = checkBit(pte, PRESENT_BIT)
	pageinfo.IsSwapped = checkBit(pte, SWAP_BIT)
	pageinfo.IsFileOrAnon = checkBit(pte, MAP_BIT)
	pageinfo.WriteProt = checkBit(pte, WRITE_PROT_BIT)
	pageinfo.MapExcl = checkBit(pte, EXCL_MAP_BIT)
	pageinfo.SoftDirty = checkBit(pte, SOFT_DIRTY_BIT)
	pageinfo.Addr = uint(int(pte&PADDR_MASK) * os.Getpagesize())

	return pageinfo
}
