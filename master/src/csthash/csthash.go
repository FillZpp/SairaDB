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


package csthash

import (
	"strconv"
	"math"
	"crypto/md5"
	"fmt"

	"config"
)

var (
	Mod uint64 = uint64(math.Pow(2, 32))
	VNodeNum uint64
	VNodeHashs []uint64
)

func Init() {
	VNodeNum, _ = strconv.ParseUint(config.ConfMap["vnode-num"], 0, 0)
	VNodeHashs = make([]uint64, 0, VNodeNum)

	tmp := Mod / VNodeNum
	var i uint64
	for i = 0; i < VNodeNum; i++ {
		VNodeHashs = append(VNodeHashs, (i + 1) * tmp)
	}
}

func FindVNode(key string) uint64 {
	b := md5.Sum([]byte(key))
	h, _ := strconv.ParseUint(fmt.Sprintf("0x%x", b[:4]), 0, 0)

	var i uint64
	for i = 0; i < VNodeNum; i++ {
		if h < VNodeHashs[i] {
			if i == 0 {
				i = VNodeNum
			}
			break
		}
	}

	return i
}


