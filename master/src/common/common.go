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

// Package common for some common use function.
package common

import (
	"os"
	"encoding/gob"
	"bytes"
	"time"
)

func IsDirExist(dir string) bool {
	flInfo, err := os.Stat(dir)

	if err != nil {
		return os.IsExist(err)
	}

	return flInfo.IsDir()
}

func DeepCopy(a, b interface{}) error {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	dec := gob.NewDecoder(buf)
	err := enc.Encode(a)
	if err != nil {
		return err
	}
	return dec.Decode(b)
}

func SetTimeout(ch chan bool, n uint) {
	time.Sleep(time.Duration(n) * time.Millisecond)
	ch<- true
}

