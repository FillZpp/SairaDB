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


package masterctl

import (
	"net"
	"fmt"
	"sync/atomic"
	
	"common"
	"slog"
)

func recvLog(ip, reason string) {
	slog.LogChan<-
		fmt.Sprintf("master controller receive task (%v): %v", ip, reason)
}

func receiveTask(idx int, ip string, connChan chan net.Conn,
	recvChan chan RecvRegister) {
	var conn net.Conn
	var err error
	var msg string
	buf := make([]byte, 1000)
	for {
		if conn == nil {
			conn = <-connChan
			msg, err = common.ConnRead(buf, conn, 500)
			if err != nil {
				recvLog(ip, err.Error())
				conn.Close()
				conn = nil
				continue
			}

			if msg != cookie {
				recvLog(ip, "wrong cookie")
				common.ConnWriteString("cookie wrong", conn, 500)
				conn.Close()
				conn = nil
				continue
			}

			err = common.ConnWriteString("ok", conn, 500)
			if err != nil {
				recvLog(ip, err.Error())
				conn.Close()
				conn = nil
				continue
			}
			atomic.AddInt32(&(MasterList[idx].Status), 1)
			recvLog(ip, "receive connected")
		}  // if conn == nil

		// TODO
	}
}


