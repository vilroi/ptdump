package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
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

type MapEntry struct {
	StartAddr uint
	EndAddr   uint
	Perms     string
	Offset    uint
	DevMajor  int
	DevMinor  int
	Inode     uint
	Path      string
}

type PageInfo struct {
	Present      bool
	IsSwapped    bool
	IsFileOrAnon bool
	WriteProt    bool
	MapExcl      bool
	SoftDirty    bool
	Addr         uint
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s pid\n", os.Args[0])
		os.Exit(1)
	}

	pid, err := strconv.Atoi(os.Args[1])
	check(err)

	maps := getMaps(pid)
	for _, m := range maps {
		pageinfo := getPageInfo(pid, m.StartAddr)
		fmt.Printf("%+v\n", m)
		fmt.Printf("%+v\n\n", pageinfo)
	}
}

// TODO: change from hardcoded value to something arch dependent
// (8 bytes assumes 64 bit
func getPageInfo(pid int, virt_addr uint) PageInfo {
	path := fmt.Sprintf("/proc/%d/pagemap", pid)
	f, err := os.Open(path)
	check(err)

	page := (virt_addr & ^(uint(0x1000 - 1))) / 0x1000
	offset := int64(page * 8)
	_, err = f.Seek(offset, os.SEEK_SET)
	check(err)

	buf := make([]byte, 8)
	_, err = f.Read(buf)
	check(err)

	pte := binary.LittleEndian.Uint64(buf)
	//fmt.Printf("0b%b\n", info)

	return newPageInfo(pte)
}

func newPageInfo(pte uint64) {
	var pageinfo PageInfo
	pageinfo.Present = checkBit(pte, PRESENT_BIT)
	pageinfo.IsSwapped = checkBit(pte, SWAP_BIT)
	pageinfo.IsFileOrAnon = checkBit(pte, MAP_BIT)
	pageinfo.WriteProt = checkBit(pte, WRITE_PROT_BIT)
	pageinfo.MapExcl = checkBit(pte, EXCL_MAP_BIT)
	pageinfo.SoftDirty = checkBit(pte, SOFT_DIRTY_BIT)
	pageinfo.Addr = uint(info & PADDR_MASK)

	return pageinfo
}

func checkBit(val uint64, bit uint) bool {
	if (uint(val) & bit) == 0 {
		return false
	}
	return true
}

func getMaps(pid int) []MapEntry {
	path := fmt.Sprintf("/proc/%d/maps", pid)
	f, err := os.Open(path)
	check(err)

	scanner := bufio.NewScanner(f)

	var maps []MapEntry
	for scanner.Scan() {
		line := scanner.Text()

		var ent MapEntry
		_, err := fmt.Sscanf(line, "%x-%x %s %x %x:%x %d %s", &ent.StartAddr, &ent.EndAddr, &ent.Perms,
			&ent.Offset, &ent.DevMajor, &ent.DevMinor, &ent.Inode, &ent.Path)

		if err == io.EOF {
			_, err = fmt.Sscanf(line, "%x-%x %s %x %x:%x %d", &ent.StartAddr, &ent.EndAddr, &ent.Perms,
				&ent.Offset, &ent.DevMajor, &ent.DevMinor, &ent.Inode)
		}
		check(err)

		maps = append(maps, ent)
	}

	return maps
}

func btoi(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
