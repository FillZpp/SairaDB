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
	"container/list"
)

type Time struct {
	Year  int
	Month int
	Day   int

	Hour   int
	Minute int
	Second int32
}

func fc(i int) string {
	if i < 10 {
		return fmt.Sprintf("0%v", i)
	}
	return fmt.Sprintf("%v", i)
}

func fc32(i int32) string {
	if i < 10 {
		return fmt.Sprintf("0%v", i)
	}
	return fmt.Sprintf("%v", i)
}

func (tm *Time) Format() string {
	return fmt.Sprintf("%v-%v-%v %v:%v:%v",
		tm.Year, fc(tm.Month), fc(tm.Day),
		fc(tm.Hour), fc(tm.Minute), fc32(tm.Second))
}

func (tm *Time) Cat() string {
	return fmt.Sprintf("%v-%v-%v_%v-%v-%v",
		tm.Year, fc(tm.Month), fc(tm.Day),
		fc(tm.Hour), fc(tm.Minute), fc32(tm.Second))
}

var (
	UnixTime int64
	Tm unsafe.Pointer

	CurQueryNum uint64 = 0

	QNOne     uint64 = 0
	QNFive    uint64 = 0
	QNFifteen uint64 = 0

	SC uint64 = 0
)

func Init() {
	ot := time.Now()
	Tm = unsafe.Pointer(&Time{
		ot.Year(), int(ot.Month()), ot.Day(),
		ot.Hour(), ot.Minute(), int32(ot.Second())})
	UnixTime = ot.UnixNano()
	go task()
}

func task() {
	var i int
	var j int32
	var ot time.Time
	tm := (*Time)(Tm)
	nt := Time{ tm.Year, tm.Month, tm.Day,
		tm.Hour, tm.Minute, tm.Second }
	fiveList := list.New()
	fifList := list.New()
	var u uint64 = 0
	for i := 0; i < 5; i ++ {
		fiveList.PushFront(u)
	}
	for i := 0; i < 10; i ++ {
		fifList.PushFront(u)
	}
	for {
		time.Sleep(time.Millisecond * 20)
		ot = time.Now()

		atomic.SwapInt64(&UnixTime, ot.UnixNano())

		j = int32(ot.Second())
		if (*Time)(Tm).Second == j {
			continue
		}
		nt.Second = j

		i = ot.Minute()
		if nt.Minute == i {
			atomic.SwapInt32(&((*Time)(Tm).Second), j)
			continue
		}
		nt.Minute = i

		// record query number
		num1 := atomic.SwapUint64(&CurQueryNum, 0)
		atomic.SwapUint64(&QNOne, num1)
		
		fiveList.PushFront(num1)
		num2 := fiveList.Remove(fiveList.Back()).(uint64)
		atomic.AddUint64(&QNFive, num1 - num2)

		fifList.PushFront(num2)
		num3 := fifList.Remove(fifList.Back()).(uint64)
		atomic.AddUint64(&QNFifteen, num1 - num3)
		
		i = ot.Hour()
		if nt.Hour == i {
			atomic.SwapPointer(&Tm, unsafe.Pointer(&nt))
			continue
		}
		nt.Hour = i
		
		i = ot.Day()
		if nt.Day == i {
			atomic.SwapPointer(&Tm, unsafe.Pointer(&nt))
			continue
		}
		nt.Day = i

		i = int(ot.Month())
		if nt.Month == i {
			atomic.SwapPointer(&Tm, unsafe.Pointer(&nt))
			continue
		}
		nt.Month = i

		i = ot.Year()
		if nt.Minute == i {
			atomic.SwapPointer(&Tm, unsafe.Pointer(&nt))
			continue
		}
		nt.Year = i
		atomic.SwapPointer(&Tm, unsafe.Pointer(&nt))
	}
}


