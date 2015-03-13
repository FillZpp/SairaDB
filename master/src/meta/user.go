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
	"errors"
	"sync/atomic"

	"common"
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

type AlterUser struct {
	Ch chan error
	AlterType string
	AlterCont []string
}

var (
	userFile *os.File
	UserChan = make(chan AlterUser, 100)
)

func initUser() {
	userFile, err := os.OpenFile(path.Join(metaDir, "/user.meta"),
		os.O_RDWR | os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"\nError:\nCan not open meta file %v:\n",
			path.Join(metaDir, "/user.meta"))
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
	go alterUserTask()
}

func alterUserTask() {
	var tmp map[string]User
	for {
		if tmp == nil {
			au := <-UserChan
			if au.AlterType == "add_user" {
				handleUserAlter((*map[string]User)(Users), au)
			} else {
				common.DeepCopy((*map[string]User)(Users), &tmp)
				handleUserAlter(&tmp, au)
			}
		} else {
			ch := make(chan bool)
			go common.SetTimeout(ch, 1)
			select {
			case au := <- UserChan:
				handleUserAlter(&tmp, au)
				continue
			case <-ch:
			}
			atomic.SwapPointer(&Databases, unsafe.Pointer(&tmp))
			tmp = nil
		}
	}
}

func handleUserAlter(users *map[string]User, au AlterUser) {
	switch au.AlterType {
	case "add_user":
		_, ok := (*users)[au.AlterCont[0]]
		if ok {
			au.Ch<- errors.New("The user already exists.")
		}
		(*users)[au.AlterCont[0]] = User{
			au.AlterCont[1],
			au.AlterCont[2],
			0,
			make(map[string]uint),
			make(map[string]uint),
		}

		// TODO
		
	default:
		au.Ch<- errors.New("Undefined alter type.")
	}
	au.Ch<- nil
}

