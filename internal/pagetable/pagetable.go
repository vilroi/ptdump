package pagetable

import (
	"fmt"
)

type PageTableEntry struct {
	virtualPage  MapEntry
	physicalPage PageInfo
}

type PageTable struct {
	entries []PageTableEntry
}

func NewPageTable(pid int) PageTable {
	var pt PageTable

	maps := getMaps(pid)
	for _, m := range maps {
		pageinfo := getPageInfo(pid, m.startAddr)

		pte := PageTableEntry{virtualPage: m, physicalPage: pageinfo}
		pt.entries = append(pt.entries, pte)
	}
	return pt
}

func (pte *PageTableEntry) Unpack() (MapEntry, PageInfo) {
	return pte.virtualPage, pte.physicalPage
}

func (pt *PageTable) Entries() []PageTableEntry {
	return pt.entries
}

func (pt *PageTable) Dump() {
	fmt.Printf("%v%27v%16v%15v%15v%15v%11v\n", "virt addr", "physical addr", "size", "perms", "present", "swapped", "path")
	for _, pte := range pt.Entries() {
		m, pageinfo := pte.Unpack()
		fmt.Printf("%0#x%#20x%20d%12s%13t%16t\t\t%s\n",
			m.startAddr, pageinfo.addr, (m.endAddr - m.startAddr),
			m.perms, pageinfo.present, pageinfo.isSwapped, m.path)
	}
}
