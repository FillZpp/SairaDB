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


package stime

import (
	"time"
	"sync/atomic"
	"fmt"
	"unsafe"

	"common"
)

type Time struct {
	Year  int
	Month int
	Day   int

	Hour   int
	Minute int
	Second int
}

func (tm *Time) Format() string {
	return fmt.Sprintf("%v-%v-%v %v:%v:%v",
		tm.Year, tm.Month, tm.Day,
		tm.Hour, tm.Minute, tm.Second)
}

var (
	UnixTime int64
	tm unsafe.Pointer
)

func Init() {
	tm = unsafe.Pointer(&Time{-1, -1, -1, -1, -1, -1})
	go task()
}

func task() {
	var i int
	var ot time.Time
	var nt Time
	for {
		time.Sleep(time.Millisecond * 100)
		ot = time.Now()

		atomic.SwapInt64(&UnixTime, ot.UnixNano())

		i = ot.Second()
		if (*Time)(tm).Second == i {
			continue
		}
		
		common.DeepCopy((*Time)(tm), &nt)
		nt.Second = i

		i = ot.Minute()
		if nt.Minute == i {
			atomic.SwapPointer(&tm, unsafe.Pointer(&nt))
			continue
		}
		nt.Minute = i
		
		i = ot.Hour()
		if nt.Hour == i {
			atomic.SwapPointer(&tm, unsafe.Pointer(&nt))
			continue
		}
		nt.Hour = i
		
		i = ot.Day()
		if nt.Day == i {
			atomic.SwapPointer(&tm, unsafe.Pointer(&nt))
			continue
		}
		nt.Day = i

		i = int(ot.Month())
		if nt.Month == i {
			atomic.SwapPointer(&tm, unsafe.Pointer(&nt))
			continue
		}
		nt.Month = i

		i = ot.Year()
		if nt.Minute == i {
			atomic.SwapPointer(&tm, unsafe.Pointer(&nt))
			continue
		}
		nt.Year = i
		atomic.SwapPointer(&tm, unsafe.Pointer(&nt))
	}
}


