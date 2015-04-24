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

	"config"
	"query"
	"csthash"
)

var (
	VNodeCtls []chan []string
	VNodeDupNum []uint64
	DupNum uint64
)

func vnodeInit() {
	DupNum, _ = strconv.ParseUint(config.ConfMap["dup-num"], 0, 0)
	VNodeCtls = make([]chan []string, csthash.VNodeNum)
	VNodeDupNum = make([]uint64, csthash.VNodeNum)
	var i uint64
	for i = 0; i < csthash.VNodeNum; i++ {
		ch := make(chan []string, 100)
		VNodeCtls[i] = ch
		VNodeDupNum[i] = 0
		go vnodeTask(i, csthash.VNodeHashs[i].Ch, ch)
	}
}

func vnodeTask(id uint64, qryChan chan query.Query, ctlChan chan []string) {
	dups := make([]string, 0, DupNum)
	dupMaster := -1

	for {
		select {
		case ctl := <-ctlChan:
			if ctl[0] == "add" {
				dups = append(dups, ctl[1])
				if len(dups) == 1 {
					dupMaster = 0
				}
			}
			// TODO
			_ = dupMaster
			continue
		case qry := <-qryChan:
			// TODO
			_ = qry
		}
	}
}

