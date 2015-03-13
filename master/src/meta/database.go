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
	nsFile *os.File
	DBChan = make(chan AlterDB, 100)
)

func initDatabase() {
	nsFile, err := os.OpenFile(path.Join(MetaDir, "/db.meta"),
		os.O_RDWR | os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"\nError:\nCan not open meta file %v:\n",
			path.Join(MetaDir, "/db.meta"))
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}

	var databases map[string]Database
	dec := gob.NewDecoder(nsFile)
	_ = dec.Decode(&databases)
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
						"key": "value",
					},
				},
			},
		}
	}

	Databases = unsafe.Pointer(&databases)
	go alterDBTask()
}

func alterDBTask() {
	var tmp map[string]Database
	for {
		if tmp == nil {
			ad := <-DBChan
			if ad.AlterType == "add_column" ||
				ad.AlterType == "add_table" ||
				ad.AlterType == "add_db" {
				handleDBAlter((*map[string]Database)(Databases), ad)
			} else {
				common.DeepCopy((*map[string]Database)(Databases),
					&tmp)
				handleDBAlter(&tmp, ad)
			}
		} else {
			ch := make(chan bool)
			go common.SetTimeout(ch, 1)
			select {
			case ad := <-DBChan:
				handleDBAlter(&tmp, ad)
				continue
			case <-ch:
			}
			atomic.SwapPointer(&Databases, unsafe.Pointer(&tmp))
			tmp = nil
		}
	}
}

func handleDBAlter(dbs *map[string]Database, ad AlterDB) {
	switch ad.AlterType {
	case "add_db":
		_, ok := (*dbs)[ad.AlterCont[0]]
		if ok {
			ad.Ch<- errors.New("The database already exists.")
		}
		(*dbs)[ad.AlterCont[0]] = Database{ ad.AlterCont[1], nil }

		// TODO
		
	default:
		ad.Ch<- errors.New("Undefined alter type.")
	}
	ad.Ch<- nil
}

	

