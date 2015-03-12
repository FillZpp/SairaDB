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
	"sync"
	"errors"

	//"common"
)

var (
	nsFile *os.File
	lock = sync.Mutex{}
)

type Table struct {
	Name string
	Key string
	Column map[string]string
}

type NameSpace struct {
	Name string
	Tables map[string]Table
}

func initNS() {
	nsFile, err := os.OpenFile(path.Join(MetaDir, "/ns.meta"),
		os.O_RDWR | os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"\nError:\nCan not open meta file %v:\n",
			path.Join(MetaDir, "/ns.meta"))
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}

	var _nameSpaces map[string]NameSpace
	dec := gob.NewDecoder(nsFile)
	_ = dec.Decode(&_nameSpaces)
	if _nameSpaces == nil {
		_nameSpaces = make(map[string]NameSpace)
	}

	_, ok := _nameSpaces["default"]
	if !ok {
		_nameSpaces["default"] = NameSpace {
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

	NameSpaces = unsafe.Pointer(&_nameSpaces)
}

func CreateNameSpace(name string) error {
	lock.Lock()
	defer lock.Unlock()
	tmp := (*map[string]NameSpace)(NameSpaces)
	_, ok := (*tmp)[name]
	if ok {
		return errors.New("The NameSpace already exists.")
	}
	(*tmp)[name] = NameSpace{ name, nil }
	return nil
}

