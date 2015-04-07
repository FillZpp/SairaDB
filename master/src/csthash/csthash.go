// SairaDB - A distributed database
// Copyright (C) 2015 by Siyu Wang
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
//	This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.


package csthash

import (
	"strconv"
	"math"
	//"fmt"

	"config"
)

type VNode struct {
	id uint64
	hash uint64
	dups []string
	dup_master string
}

var (
	Mod uint64 = uint64(math.Pow(2, 32))
	VNodeNum uint64
	DupNum uint64

	VNodes []VNode
)

func Init() {
	VNodeNum, _ = strconv.ParseUint(config.ConfMap["vnode-num"], 0, 0)
	DupNum, _ = strconv.ParseUint(config.ConfMap["dup-num"], 0, 0)
	VNodes = make([]VNode, 0, VNodeNum)

	tmp := Mod / VNodeNum
	var i uint64
	for i = 0; i < VNodeNum; i++ {
		VNodes = append(VNodes, VNode{
			i,
			i * tmp,
			make([]string, 0, DupNum),
			"",
		})
	}
}


