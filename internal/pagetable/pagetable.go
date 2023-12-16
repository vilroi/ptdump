package pagetable

import (
	"fmt"
)

type PageTableEntry struct {
	VirtualPage  MapEntry
	PhysicalPage PageInfo
}

type PageTable struct {
	entries []PageTableEntry
}

func (pte *PageTableEntry) Unpack() (MapEntry, PageInfo) {
	return pte.VirtualPage, pte.PhysicalPage
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
