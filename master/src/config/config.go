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


// Package config generate CFGMap from configure file and default configures.
package config

import (
	"fmt"
	"sync"
)

var once sync.Once
var CFGMap = make(map[string]string)
var MasterList = make([]string, 0, 10)

// Initialize CFGMap for only once
func Init() {
	once.Do(func() {
		fmt.Println("Config init")

		CFGMap["etc_dir"] = Prefix

		readMasters()
		fmt.Println(MasterList)
		
		fileMap := readConfFile()
		fmt.Println(fileMap)

	})
}

