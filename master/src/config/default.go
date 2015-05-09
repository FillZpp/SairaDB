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


package config

import (
	"path"
	"crypto/md5"
	"fmt"
)

var BoolConfs = []string {
	"serialize",
	"local",
	"delete-old-log",
}

var UintConfs = []string {
	"port-client",
	"port-master",
	"port-slave",
	"vnode-num",
	"dup-num",
}

func defaultConf() {
	ConfMap["serialize"] = "on"
	ConfMap["local"] = "off"
	ConfMap["delete-old-log"] = "on"
	
	ConfMap["data-dir"] = path.Join(HomeDir, "/saira_data")
	ConfMap["log-level"] = "error"

	ConfMap["port-client"] = "4400"
	ConfMap["port-master"] = "4401"
	ConfMap["port-slave"]  = "4402"

	ckMd5 := fmt.Sprintf("%x", md5.Sum([]byte("")))
	ConfMap["client-cookie"] = ckMd5
	ConfMap["master-cookie"] = ckMd5
	ConfMap["slave-cookie"] = ckMd5

	ConfMap["dup-num"] = "3"
	ConfMap["vnode-num"] = "16"

	// TODO
	// readDeadLine
}

