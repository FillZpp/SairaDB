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


// Package meta for meta data control.
package meta

import (
	"fmt"
	"os"
	"path"
	"unsafe"

	"common"
	"config"
)

var (
	metaDir string
	
	Databases unsafe.Pointer  // map[string]NameSpace
	Users unsafe.Pointer      // map[string]User

	ToClose = make(chan bool, 2)
	GetEnd = make(chan bool)
)

func Init() {
	metaDir = path.Join(config.ConfMap["data-dir"], "/master/meta/")

	if !common.IsDirExist(metaDir) {
		err := os.MkdirAll(metaDir, 0700)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"\nError:\nCan not create meta dir %v:\n%v\n",
				metaDir,
				err.Error())
			os.Exit(3)
		}
	}

	initDatabase()
	initUser()

}


