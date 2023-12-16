package pagetable

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var GranularFlag = false

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

func (m *MapEntry) Size() int {
	return int(m.EndAddr - m.StartAddr)
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
		tmp.StartAddr += uint(i * pagesize)
		tmp.EndAddr = tmp.StartAddr + uint(pagesize)

		maps[i] = tmp
	}

	return maps
}
