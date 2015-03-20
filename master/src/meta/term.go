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
	"fmt"
	"os"
	"path"
	"strconv"
	"sync/atomic"

	"common"
)

var (
	termFile *os.File
	TermChan = make(chan bool, 100)
	closeTerm = make(chan bool, 1)
	termClosed = make(chan bool, 1)
)

func initTerm() {
	var err error
	termFile, err = os.OpenFile(path.Join(metaDir, "/term.meta"),
		os.O_RDWR | os.O_CREATE, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"\nError:\nCan not open meta file %v:\n%v\n",
			path.Join(metaDir, "/term.meta"),
			err.Error())
		os.Exit(3)
	}

	bt := make([]byte, 40)
	n, err := termFile.Read(bt)
	if err == nil && n > 0 {
		s := string(bt[:n])
		Term, _ = strconv.ParseUint(s, 0, 0)
	}

	go termTask()
}

func syncTermFile() {
	termFile.Seek(0, 0)
	termFile.WriteString(fmt.Sprintf("%v", atomic.LoadUint64(&Term)))
}

func termTask() {
	status := true
	for {
		if status {
			select {
			case <-TermChan:
				status = false
			case <-closeTerm:
				termFile.Close()
				termClosed<- true
				for {
					<-TermChan
				}
			}
		} else {
			ch := make(chan bool, 1)
			go common.SetTimeout(ch, 10)
			select {
			case <-TermChan:
				continue
			case <-ch:
			}

			syncTermFile()
			status = true
		}
	}
}
	
