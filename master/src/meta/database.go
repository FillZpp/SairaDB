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

type Table struct {
	Name string
	Key string
	Column map[string]string
}

type Database struct {
	Name string
	Tables map[string]Table
}

type AlterDB struct {
	Ch chan error
	AlterType string
	AlterCont []string
}

var (
	dbFile *os.File
	DBChan = make(chan AlterDB, 100)
)

func initDatabase() {
	var err error
	dbFile, err = os.OpenFile(path.Join(metaDir, "/db.meta"),
		os.O_RDWR | os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"\nError:\nCan not open meta file %v:\n%v\n",
			path.Join(metaDir, "/db.meta"),
			err.Error())
		os.Exit(3)
	}

	var databases map[string]Database
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
		buf := bytes.NewBufferString(strs[1])
		
		dec := gob.NewDecoder(buf)
		_ = dec.Decode(&databases)
	}
	
	if databases == nil {
		databases = make(map[string]Database)
	}

	_, ok := databases["default"]
	if !ok {
		databases["default"] = Database{
			"default",
			map[string]Table {
				"kv": Table {
					"kv",
					"key",
					map[string]string {
						"key": "Any",
					},
				},
			},
		}
	}

	Databases = unsafe.Pointer(&databases)
	go alterDBTask()
}

func syncDBFile() {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode((*map[string]Database)(atomic.LoadPointer(&Databases)))
	s := buf.String()
	s = strconv.Itoa(len(s)) + ";" + s
	
	dbFile.Seek(0, 0)
	dbFile.WriteString(s)
}

func alterDBTask() {
	syncDBFile()
	var tmp map[string]Database
	var ad AlterDB
	for {
		if tmp == nil {
			select {
			case <-ToClose:
				dbFile.Close()
				GotIt<- true
				for {
					ad = <-DBChan
					ad.Ch<- errors.New("This master is to close.")
				}
			case ad = <-DBChan:
			}
			if ad.AlterType == "add_column" ||
				ad.AlterType == "add_table" ||
				ad.AlterType == "add_db" {
				handleDBAlter((*map[string]Database)(Databases), ad)
			} else {
				common.DeepCopy((*map[string]Database)(Databases),
					&tmp)
				if !handleDBAlter(&tmp, ad) {
					tmp = nil
				}
			}
		} else {
			ch := make(chan bool)
			go common.SetTimeout(ch, 1)
			select {
			case ad = <-DBChan:
				handleDBAlter(&tmp, ad)
				continue
			case <-ch:
			}
			atomic.StorePointer(&Databases, unsafe.Pointer(&tmp))
			tmp = nil
			syncDBFile()
		}
	}
}

func handleDBAlter(dbs *map[string]Database, ad AlterDB) bool {
	switch ad.AlterType {
	case "add_db":
		_, ok := (*dbs)[ad.AlterCont[0]]
		if ok {
			ad.Ch<- errors.New("The database already exists.")
			return false
		}
		(*dbs)[ad.AlterCont[0]] = Database{ ad.AlterCont[1], nil }

		// TODO
		
	default:
		ad.Ch<- errors.New("Undefined alter type.")
		return false
	}
	ad.Ch<- nil
	return true
}

	

