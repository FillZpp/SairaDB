// SairaDB - A distributed database
// Copyright (C) 2015 by Siyu Wang
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.


package slavectl

import (
	"net"
	"fmt"
	"time"
	"sync"
	"sync/atomic"
	"encoding/json"

	"slog"
	"common"
	"query"
	"csthash"
)

func handlerLog(ip, reason string) {
	slog.LogChan<-
		fmt.Sprintf("slave controller handle (%v): %v", ip, reason)
}

func slaveHandler(conn net.Conn) {
	defer conn.Close()
	var msg string
	var status string
	buf := make([]byte, 1000)
	ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		slog.LogChan<- "slave controller address split error: " + err.Error()
		return
	}

	msg, err = common.ConnRead(buf, conn, 1000)
	if err != nil {
		handlerLog(ip, err.Error())
		return
	}

	if msg != cookie {
		handlerLog(ip, "wrong cookie")
		common.ConnWriteString("wrong cookie", conn, 1000)
		return
	}

	// check
	mutex.Lock()
	slv, ok := Slaves[ip]
	if !ok {
		var rwMutex sync.RWMutex
		slv = &Slave{
			ip,
			make([]uint64, 0),
			rwMutex,
			make(chan query.Query, 10000),
			make(chan RecvRegister, 10000),
			0,
			0,
		}
		Slaves[ip] = slv
	}
	mutex.Unlock()
	
	if atomic.CompareAndSwapInt32(&slv.sendStatus, 0, 1) {
		status = "send"
	} else if atomic.CompareAndSwapInt32(&slv.recvStatus, 0, 1) {
		status = "recv"
	} else {
		status = "no need"
	}

	err = common.ConnWriteString(status, conn, 1000)
	if err != nil {
		handlerLog(ip, err.Error())
		return
	} else if status == "no need" {
		handlerLog(ip, "no need")
		return
	}
	
	handlerLog(ip, "connected")

	time.Sleep(time.Hour)
	if status == "send" {
		sendSlave(conn, slv)
	} else {
		recvSlave(conn, slv)
	}
}

func sendSlave(conn net.Conn, slv *Slave) {
	
}

func recvSlave(conn net.Conn, slv *Slave) {
	buf := make([]byte, 1000)
	msg, err := common.ConnRead(buf, conn, 1000)
	if err != nil {
		handlerLog(slv.ip, err.Error())
		atomic.StoreInt32(&slv.recvStatus, 0)
	}
	
	var arr []uint64
	err = json.Unmarshal([]byte(msg), &arr)
	if err == nil {
		vnodes := make([]uint64, 0, len(arr))
		for _, i := range arr {
			if i < csthash.VNodeNum {
				vnodes = append(vnodes, i)
			}
		}
		slv.rwMutex.Lock()
		slv.vnodes = vnodes
		slv.rwMutex.Unlock()
	}

	err = common.ConnWriteString("ok", conn, 1000)
}


