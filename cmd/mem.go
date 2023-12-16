package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
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

type PageTableEntry struct {
	VirtualPage  MapEntry
	PhysicalPage PageInfo
}

func (pte *PageTableEntry) Unpack() (MapEntry, PageInfo) {
	return pte.VirtualPage, pte.PhysicalPage
}

type PageTable struct {
	entries []PageTableEntry
}

func (m *MapEntry) Size() int {
	return int(m.EndAddr - m.StartAddr)
}

func NewPageTable(pid int) PageTable {
	var pt PageTable

	maps := getMaps(pid)
	for _, m := range maps {
		pageinfo := getPageInfo(pid, m.StartAddr)

		pte := PageTableEntry{VirtualPage: m, PhysicalPage: pageinfo}
		pt.entries = append(pt.entries, pte)
	}

	return pt
}

func (pt *PageTable) Entries() []PageTableEntry {
	return pt.entries
}

func (pt *PageTable) Dump() {
	fmt.Printf("%v%30v%13v%15v%15v%15v%11v\n", "virt addr", "physical addr", "size", "perms", "present", "swapped", "path")
	for _, pte := range pt.Entries() {
		m, pageinfo := pte.Unpack()
		fmt.Printf("%0#x%#20x%20d%12s%13t%16t\t\t%s\n",
			m.StartAddr, pageinfo.Addr, (m.EndAddr - m.StartAddr),
			m.Perms, pageinfo.Present, pageinfo.IsSwapped, m.Path)
	}
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

		/*TODO: Ugly. Clean up and make simpler */
		if granular_flag {
			pagesize := os.Getpagesize()
			if ent.Size() > pagesize {
				split := splitPages(ent)
				maps = append(maps, split...)
			} else {
				maps = append(maps, ent)
			}
		} else {
			maps = append(maps, ent)
		}
	}

	return maps
}

func splitPages(m MapEntry) []MapEntry {
	pagesize := os.Getpagesize()
	count := m.Size() / pagesize

	maps := make([]MapEntry, count)
	for i := 0; i < count; i++ {
		tmp := m
		tmp.StartAddr += uint(i * pagesize)
		tmp.EndAddr = tmp.StartAddr + uint(pagesize)

		maps[i] = tmp
	}

	return maps
}
