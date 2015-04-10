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


package slavectl

import (
	"net"
	"fmt"
	"time"
	"sync/atomic"

	"slog"
	"common"
	"query"
)

func handlerLog(ip, reason string) {
	slog.LogChan<-
		fmt.Sprintf("slave controller handle (%v): %v", ip, reason)
}

func slaveHandler(conn net.Conn) {
	defer conn.Close()
	var err error
	var msg string
	var ip string
	var status string
	buf := make([]byte, 1000)
	ip, _, err = net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		slog.LogChan<- "slave controller address split error: " + err.Error()
		return
	}

	// check
	mutex.Lock()
	slv, ok := Slaves[ip]
	if !ok {
		slv = &Slave{
			ip,
			make([]uint64, 0),
			make(chan query.Query, 10000),
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

	msg, err = common.ConnRead(buf, conn, 100)
	if err != nil {
		handlerLog(ip, err.Error())
		return
	}
	
	if msg != cookie {
		handlerLog(ip, "wrong cookie")
		common.ConnWriteString("wrong cookie", conn, 100)
		return
	}
		
	err = common.ConnWriteString(status, conn, 100)
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
		sendSlave(conn)
	} else {
		recvSlave(conn)
	}
}

func sendSlave(conn net.Conn) {
	
}

func recvSlave(conn net.Conn) {
	
}


