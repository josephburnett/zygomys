// Inferno utils/5l/obj.c
// http://code.google.com/p/inferno-os/source/browse/utils/5l/obj.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package s390x

import (
	"cmd/avail/obj"
	"cmd/avail/sys"
	"cmd/link/avail/ld"
	"fmt"
)

// Reading object files.

func Main() {
	linkarchinit()
	ld.Ldmain()
}

func linkarchinit() {
	ld.SysArch = sys.ArchS390X

	ld.Thearch.Funcalign = FuncAlign
	ld.Thearch.Maxalign = MaxAlign
	ld.Thearch.Minalign = MinAlign
	ld.Thearch.Dwarfregsp = DWARFREGSP
	ld.Thearch.Dwarfreglr = DWARFREGLR

	ld.Thearch.Adddynrel = adddynrel
	ld.Thearch.Archinit = archinit
	ld.Thearch.Archreloc = archreloc
	ld.Thearch.Archrelocvariant = archrelocvariant
	ld.Thearch.Asmb = asmb // in asm.go
	ld.Thearch.Elfreloc1 = elfreloc1
	ld.Thearch.Elfsetupplt = elfsetupplt
	ld.Thearch.Gentext = gentext
	ld.Thearch.Machoreloc1 = machoreloc1
	ld.Thearch.Lput = ld.Lputb
	ld.Thearch.Wput = ld.Wputb
	ld.Thearch.Vput = ld.Vputb
	ld.Thearch.Append16 = ld.Append16b
	ld.Thearch.Append32 = ld.Append32b
	ld.Thearch.Append64 = ld.Append64b

	ld.Thearch.Linuxdynld = "/lib64/ld64.so.1"

	// not relevant for s390x
	ld.Thearch.Freebsddynld = "XXX"
	ld.Thearch.Openbsddynld = "XXX"
	ld.Thearch.Netbsddynld = "XXX"
	ld.Thearch.Dragonflydynld = "XXX"
	ld.Thearch.Solarisdynld = "XXX"
}

func archinit() {
	// getgoextlinkenabled is based on GO_EXTLINK_ENABLED when
	// Go was built; see ../../make.bash.
	if ld.Linkmode == ld.LinkAuto && obj.Getgoextlinkenabled() == "0" {
		ld.Linkmode = ld.LinkInternal
	}

	if ld.Buildmode == ld.BuildmodeCArchive || ld.Buildmode == ld.BuildmodeCShared || ld.DynlinkingGo() {
		ld.Linkmode = ld.LinkExternal
	}

	switch ld.HEADTYPE {
	default:
		ld.Exitf("unknown -H option: %v", ld.HEADTYPE)

	case obj.Hlinux: // s390x ELF
		ld.Elfinit()
		ld.HEADR = ld.ELFRESERVE
		if ld.INITTEXT == -1 {
			ld.INITTEXT = 0x10000 + int64(ld.HEADR)
		}
		if ld.INITDAT == -1 {
			ld.INITDAT = 0
		}
		if ld.INITRND == -1 {
			ld.INITRND = 0x10000
		}
	}

	if ld.INITDAT != 0 && ld.INITRND != 0 {
		fmt.Printf("warning: -D0x%x is ignored because of -R0x%x\n", uint64(ld.INITDAT), uint32(ld.INITRND))
	}
}
