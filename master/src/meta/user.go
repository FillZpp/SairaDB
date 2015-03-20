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
	"strings"
	"strconv"
	"bytes"

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
	var err error
	userFile, err = os.OpenFile(path.Join(metaDir, "/user.meta"),
		os.O_RDWR | os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"\nError:\nCan not open meta file %v:\n%v\n",
			path.Join(metaDir, "/user.meta"),
			err.Error())
		os.Exit(3)
	}

	var meta []string
	var users map[string]User
	bt := make([]byte, 10)
	n, err := userFile.Read(bt)
	if err == nil && n > 0 {
		strs := strings.SplitN(string(bt), ";", 2)
		length, _ := strconv.Atoi(strs[0])
		if length + len(strs[0]) + 1 > 10 {
			bt = make([]byte, length + len(strs[0]) + 1 - 10)
			userFile.Read(bt)
			strs[1] += string(bt)
		}
		
		metaBuf := bytes.NewBufferString(strs[1])
		metaDec := gob.NewDecoder(metaBuf)
		_ = metaDec.Decode(&meta)

		term, _ := strconv.Atoi(meta[0])
		if int32(term) > Term {
			Term = int32(term)
		}

		userBuf := bytes.NewBufferString(meta[1])
		userDec := gob.NewDecoder(userBuf)
		_ = userDec.Decode(&users)
	}
	
	if users == nil {
		users = make(map[string]User)
	}

	_, ok := users["root"]
	if !ok {
		users["root"] = User {
			"root",
			"",
			Super,
			map[string]uint{},
			map[string]uint{},
		}
	}

	Users = unsafe.Pointer(&users)
	go alterUserTask()
}

func syncUserFile() {
	userBuf := new(bytes.Buffer)
	userEnc := gob.NewEncoder(userBuf)
	userEnc.Encode((*map[string]User)(Users))
	userStr := userBuf.String()

	metaBuf := new(bytes.Buffer)
	metaEnc := gob.NewEncoder(metaBuf)
	metaEnc.Encode([]string{
		strconv.Itoa(int(atomic.LoadInt32(&Term))),
		userStr,
	})

	s := metaBuf.String()
	s = strconv.Itoa(len(s)) + ";" + s

	userFile.Seek(0, 0)
	userFile.WriteString(s)
}

func alterUserTask() {
	syncUserFile()
	var tmp *map[string]User
	var au AlterUser
	for {
		if tmp == nil {
			select {
			case <-ToClose:
				userFile.Close()
				GetEnd<- true
				for {
					au := <-UserChan
					au.Ch<- errors.New("This master is to close.")
				}
			case au = <-UserChan:
			}
			var copy map[string]User
			common.DeepCopy((*map[string]User)(Users), &copy)
			tmp = &copy
			if !handleUserAlter(tmp, au) {
				tmp = nil
			}
		} else {
			ch := make(chan bool)
			go common.SetTimeout(ch, 10)
			select {
			case au = <- UserChan:
				handleUserAlter(tmp, au)
				continue
			case <-ch:
			}
			
			atomic.StorePointer(&Users, unsafe.Pointer(tmp))
			atomic.AddInt32(&Term, 1)
			tmp = nil
			syncUserFile()
		}
	}
}

func handleUserAlter(users *map[string]User, au AlterUser) bool {
	switch au.AlterType {
	case "add_user":
		_, ok := (*users)[au.AlterCont[0]]
		if ok {
			au.Ch<- errors.New("The user already exists.")
			return false
		}
		(*users)[au.AlterCont[0]] = User{
			au.AlterCont[0],
			au.AlterCont[1],
			0,
			make(map[string]uint),
			make(map[string]uint),
		}

		// TODO
		
	default:
		au.Ch<- errors.New("Undefined alter type.")
		return false
	}
	au.Ch<- nil
	return true
}

