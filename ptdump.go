package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
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

type PageType int

type PageInfo struct {
	Present      bool
	IsSwapped    bool
	IsFileOrAnon bool
	WriteProt    bool
	MapExcl      bool
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

	info := binary.LittleEndian.Uint64(buf)
	//fmt.Printf("0b%b\n", info)

	var pageinfo PageInfo
	pageinfo.Present = btoi(int(info >> 63))
	pageinfo.IsSwapped = btoi(int(info>>62) & 1)
	pageinfo.IsFileOrAnon = btoi(int(info>>61) & 1)
	pageinfo.WriteProt = btoi(int(info>>57) & 1)
	pageinfo.MapExcl = btoi(int(info>>56) & 1)
	pageinfo.Addr = uint(info & ((1 << 56) - 1))

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
