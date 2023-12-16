package pagetable

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var GranularFlag = false

type MapEntry struct {
	startAddr uint
	endAddr   uint
	perms     string
	offset    uint
	devMajor  int
	devMinor  int
	inode     uint
	path      string
}

func (m *MapEntry) Size() int {
	return int(m.endAddr - m.startAddr)
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
		_, err := fmt.Sscanf(line, "%x-%x %s %x %x:%x %d %s", &ent.startAddr, &ent.endAddr, &ent.perms,
			&ent.offset, &ent.devMajor, &ent.devMinor, &ent.inode, &ent.path)

		if err == io.EOF {
			_, err = fmt.Sscanf(line, "%x-%x %s %x %x:%x %d", &ent.startAddr, &ent.endAddr, &ent.perms,
				&ent.offset, &ent.devMajor, &ent.devMinor, &ent.inode)
		}
		check(err)

		/*TODO: Ugly. Clean up and make simpler */
		if GranularFlag {
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
		tmp.startAddr += uint(i * pagesize)
		tmp.endAddr = tmp.startAddr + uint(pagesize)

		maps[i] = tmp
	}

	return maps
}
