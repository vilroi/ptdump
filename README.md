# ptdump


## Summary
ptdump is a tool to dump information about a process's address space on linux.

The information dumped includes the following:

- virtual address of each mapped region in the process's address space
- physical address of corresponding physical page frame
- size of mapped region
- permissions
- if the page is present in memory
- if it has been swapped out
- path

Since this tool is dependent on /proc/[pid]/maps and /proc/[pid]/pagemap, it may not be protable to other \*nix environments.

Also, the tool must be run with root privileges, or it will not display the physical addresses.

## Usage

By default, ptdump a single pid as an argument

```console
$ sudo ./ptdump 89088
virt addr              physical addr            size          perms        present        swapped 	path
0x562dae786000         0x1b26c2000              192512        r--p         true           false		/usr/bin/bash
0x562dae7b5000         0x12dc67000              790528        r-xp         true           false		/usr/bin/bash
0x562dae876000         0x1b25ea000              229376        r--p         true           false		/usr/bin/bash
0x562dae8ae000         0x11c87f000               16384        r--p         true           false		/usr/bin/bash
0x562dae8b2000         0x1ec200000               36864        rw-p         true           false		/usr/bin/bash
0x562dae8bb000         0x21d296000               45056        rw-p         true           false		
0x562dafaa7000         0x1ec3ec000             1638400        rw-p         true           false		[heap]
0x7ff7dc200000         0x187284000             3051520        r--p         true           false		/usr/lib/locale/locale-archive
0x7ff7dc683000         0x22a93b000               12288        rw-p         true           false		
0x7ff7dc686000         0x13b62a000              155648        r--p         true           false		/usr/lib/x86_64-linux-gnu/libc.so.6
0x7ff7dc6ac000         0x3bee76000             1396736        r-xp         true           false		/usr/lib/x86_64-linux-gnu/libc.so.6
0x7ff7dc801000         0x1982d1000              339968        r--p         true           false		/usr/lib/x86_64-linux-gnu/libc.so.6
0x7ff7dc854000         0x244dee000               16384        r--p         true           false		/usr/lib/x86_64-linux-gnu/libc.so.6
0x7ff7dc858000         0x136fae000                8192        rw-p         true           false		/usr/lib/x86_64-linux-gnu/libc.so.6
0x7ff7dc85a000         0x24ab9d000               53248        rw-p         true           false		
0x7ff7dc867000         0x3af009000               61440        r--p         true           false		/usr/lib/x86_64-linux-gnu/libtinfo.so.6.4
0x7ff7dc876000         0x153892000               69632        r-xp         true           false		/usr/lib/x86_64-linux-gnu/libtinfo.so.6.4
0x7ff7dc887000         0x1873f9000               57344        r--p         true           false		/usr/lib/x86_64-linux-gnu/libtinfo.so.6.4
0x7ff7dc895000         0x204482000               16384        r--p         true           false		/usr/lib/x86_64-linux-gnu/libtinfo.so.6.4
0x7ff7dc899000         0x1be3ab000                4096        rw-p         true           false		/usr/lib/x86_64-linux-gnu/libtinfo.so.6.4
0x7ff7dc8a7000         0x1440c7000               28672        r--s         true           false		/usr/lib/x86_64-linux-gnu/gconv/gconv-modules.cache
0x7ff7dc8ae000         0x1e7368000                8192        rw-p         true           false		
0x7ff7dc8b0000         0x3bee36000                4096        r--p         true           false		/usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
0x7ff7dc8b1000         0x10ace9000              151552        r-xp         true           false		/usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
0x7ff7dc8d6000         0x153938000               40960        r--p         true           false		/usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
0x7ff7dc8e0000         0x1c5673000                8192        r--p         true           false		/usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
0x7ff7dc8e2000         0x244b53000                8192        rw-p         true           false		/usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
0x7fff4b9bb000                 0x0              135168        rw-p        false           false		[stack]
0x7fff4b9ee000                 0x0               16384        r--p        false           false		[vvar]
0x7fff4b9f2000         0x330a0b000                8192        r-xp         true           false		[vdso]
```

The `-g` option flag can be set to show every single page (pages are grouped together by default). 

Warning: This option will spit out a lot of lines.
``` console
$ sudo ./ptdump -g 23424 | head 
virt addr              physical addr              size        perms        present        swapped   	  path
0x562dae786000         0x1b26c2000                4096        r--p         true           false		  /usr/bin/bash
0x562dae787000         0x1b26c6000                4096        r--p         true           false		  /usr/bin/bash
0x562dae788000         0x1b26c1000                4096        r--p         true           false		  /usr/bin/bash
0x562dae789000         0x1b26c7000                4096        r--p         true           false		  /usr/bin/bash
0x562dae78a000         0x1b26d2000                4096        r--p         true           false		  /usr/bin/bash
0x562dae78b000         0x1b26cd000                4096        r--p         true           false		  /usr/bin/bash
0x562dae78c000         0x1b26d7000                4096        r--p         true           false		  /usr/bin/bash
0x562dae78d000         0x1b24f2000                4096        r--p         true           false		  /usr/bin/bash
0x562dae78e000         0x129685000                4096        r--p         true           false		  /usr/bin/bash
```

Other usage includes combining it with watch(1) to watch physical page allocation (demand paging) happen in real time.
