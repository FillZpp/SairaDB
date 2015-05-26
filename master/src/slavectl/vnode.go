// SairaDB - A distributed database
// Copyright (C) 2015 by Siyu Wang
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.


package slavectl

import (
	"strconv"
	"sync/atomic"
	"unsafe"

	"config"
	"csthash"
)

type VNodeCtl struct {
	cont []string
	resChan chan string
}

var (
	DupNum uint64
	VNodeCtls []chan VNodeCtl
	VNodeDupNum []uint64
	VNodeDupMaster []unsafe.Pointer // string
)

func vnodeInit() {
	DupNum, _ = strconv.ParseUint(config.ConfMap["dup-num"], 0, 0)
	VNodeCtls = make([]chan VNodeCtl, csthash.VNodeNum)
	VNodeDupNum = make([]uint64, csthash.VNodeNum)
	VNodeDupMaster = make([]unsafe.Pointer, csthash.VNodeNum)
	var i uint64
	for i = 0; i < csthash.VNodeNum; i++ {
		ch := make(chan VNodeCtl, 100)
		VNodeCtls[i] = ch
		VNodeDupNum[i] = 0
		go vnodeTask(i, ch)
	}
}

func vnodeTask(id uint64, ctlChan chan VNodeCtl) {
	dups := make(map[string]int)
	dupMaster := ""
	var n uint64 = 0

	for {
		ctl := <-ctlChan
		switch ctl.cont[0] {
		case "add": 
			dups[ctl.cont[1]] = 1
			if len(dups) == 1 {
				dupMaster = ctl.cont[1]
			}
			n += 1
		case "del":
			if ctl.cont[1] == dupMaster {
				if len(dups) == 1 {
					dupMaster = ""
				} else {
					// TODO
					// Check other duplicate term
					// and choose new master
				}
			}
			delete(dups, ctl.cont[1])
			n -= 1
		}
		atomic.StoreUint64(&VNodeDupNum[id], n)
		ctl.resChan<- dupMaster
	}
}

