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
	"encoding/json"
	"unsafe"
	"errors"
	"sync/atomic"
	"strings"
	"strconv"

	"common"
)

type AlterDB struct {
	Ch chan error
	AlterType string
	AlterCont []string
}

var (
	dbFile *os.File
	DBChan = make(chan AlterDB, 100)
	closeDB = make(chan bool, 1)
	dbClosed = make(chan bool, 1)
)

func initDatabase() {
	var err error
	dbFile, err = os.OpenFile(path.Join(metaDir, "/db.meta"),
		os.O_RDWR | os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error:\nCan not open meta file %v:\n%v\n",
			path.Join(metaDir, "/db.meta"),
			err.Error())
		os.Exit(3)
	}

	var databases map[string]int
	bt := make([]byte, 10)
	n, err := dbFile.Read(bt)
	if err == nil && n > 0 {
		strs := strings.SplitN(string(bt), ";", 2)
		length, _ := strconv.Atoi(strs[0])
		if length + len(strs[0]) + 1 > 10 {
			bt = make([]byte, length + len(strs[0]) + 1 - 10)
			dbFile.Read(bt)
			strs[1] += string(bt)
		}

		json.Unmarshal([]byte(strs[1]), &databases)
	}
	
	if databases == nil {
		databases = make(map[string]int)
		databases["default"] = 1
	}
	
	Databases = unsafe.Pointer(&databases)
	go alterDBTask()
}

func syncDBFile() {
	b, _ := json.Marshal((*map[string]int)(atomic.LoadPointer(&Databases)))
	s := string(b)
	atomic.StorePointer(&DBEncode, unsafe.Pointer(&s))
	s2 := strconv.Itoa(len(s)) + ";" + s
	
	dbFile.Seek(0, 0)
	dbFile.WriteString(s2)
}

func alterDBTask() {
	syncDBFile()
	var tmp *map[string]int
	var ad AlterDB
	for {
		if tmp == nil {
			select {
			case <-closeDB:
				dbFile.Close()
				dbClosed<- true
				for {
					ad = <-DBChan
					ad.Ch<- errors.New("This master is to close.")
				}
			case ad = <-DBChan:
			}
			var copy map[string]int
			common.DeepCopy((*map[string]int)(Databases), &copy)
			tmp = &copy
			if !handleDBAlter(tmp, ad) {
				tmp = nil
			}
		} else {
			ch := make(chan bool, 1)
			go common.SetTimeout(ch, 10)
			select {
			case ad = <-DBChan:
				handleDBAlter(tmp, ad)
				continue
			case <-ch:
			}

			atomic.StorePointer(&Databases, unsafe.Pointer(tmp))
			tmp = nil
			syncDBFile()
		}
	}
}

func handleDBAlter(dbs *map[string]int, ad AlterDB) bool {
	switch ad.AlterType {
	case "add_db":
		_, ok := (*dbs)[ad.AlterCont[0]]
		if ok {
			ad.Ch<- errors.New("The database already exists.")
			return false
		}
		(*dbs)[ad.AlterCont[0]] = 1
		// TODO
		
	default:
		ad.Ch<- errors.New("Undefined alter type.")
		return false
	}
	ad.Ch<- nil
	return true
}

	

