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


package meta

import (
	"os"
	"path"
	"fmt"
	"encoding/gob"
	"unsafe"

	"common"
)

var (
	userFile *os.File
	UserChan = make(chan []string, 100)
)

const (
	Read  = 1
	Write = 2
	Alter = 4
	Delete = 8

	Create = 16
	Drop = 32
	
	CreateUser = 64
	Super = 128
)

type User struct {
	Name string
	Password string
	
	GlobalAuth uint
	NameSpaceAuth map[string]uint
	TableAuth map[string]uint
}

func initUser() {
	userFile, err := os.OpenFile(path.Join(MetaDir, "/user.meta"),
		os.O_RDWR | os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"\nError:\nCan not open meta file %v:\n",
			path.Join(MetaDir, "/user.meta"))
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}

	var _users map[string]User
	dec := gob.NewDecoder(userFile)
	_ = dec.Decode(&_users)
	if _users == nil {
		_users = make(map[string]User)
	}

	_, ok := _users["root"]
	if !ok {
		_users["root"] = User {
			"root",
			"",
			Super,
			map[string]uint{},
			map[string]uint{},
		}
	}

	Users = unsafe.Pointer(&_users)

	go userCoroutine()
}

func userCoroutine() {
	common.IsDirExist("/")
}

