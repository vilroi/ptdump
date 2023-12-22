package pagetable

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

var (
	GranularFlag      = false
	StartAddr    uint = 0
	Length       uint = 0
	EndAddr      uint = uint(math.Pow(2, 64)) - 1
	pagesize          = uint(os.Getpagesize())
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

		if addrWithinRange(&ent) == false {
			continue
		}

		/*TODO: Ugly. Clean up and make simpler */
		if GranularFlag {
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

func addrWithinRange(m *MapEntry) bool {
	if Length != 0 {
		EndAddr = pageAlign(StartAddr + Length)
	}
	if StartAddr <= m.startAddr && m.endAddr <= EndAddr {
		return true
	}
	return false
}
