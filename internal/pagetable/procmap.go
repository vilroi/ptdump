package pagetable

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	GranularFlag = false
	pagesize     = uint(os.Getpagesize())
)

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

func (m *MapEntry) Size() uint {
	return m.endAddr - m.startAddr
}

// TODO: this funtion is starting to get hacky...think of way to make simpler
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

		if GranularFlag {
			if ent.Size() > pagesize {
				split := splitPages(ent)
				maps = append(maps, split...)
				continue
			}
		}
		maps = append(maps, ent)

	}

	return maps
}

func splitPages(m MapEntry) []MapEntry {
	count := int(m.Size() / pagesize)

	maps := make([]MapEntry, count)
	for i := 0; i < count; i++ {
		tmp := m
		tmp.startAddr += uint(i * int(pagesize))
		tmp.endAddr = tmp.startAddr + uint(pagesize)

		maps[i] = tmp
	}

	return maps
}
