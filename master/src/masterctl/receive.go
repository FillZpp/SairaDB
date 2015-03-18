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


package masterctl

import (
	"net"
	"fmt"
	"time"
	
	"common"
	"slog"
)

func recvLog(ip, reason string) {
	slog.LogChan<-
		fmt.Sprintf("master controller receive task (%v): %v",
		ip, reason)
}

func receiveTask(ip string, ch chan net.Conn) {
	var conn net.Conn
	var err error
	var msg string
	buf := make([]byte, 1000)
	for {
		if conn == nil {
			conn = <-ch
			msg, err = common.ConnRead(buf, conn, 100)
			if err != nil {
				recvLog(ip, err.Error())
				conn.Close()
				continue
			}

			if msg != cookie {
				recvLog(ip, "wrong cookie")
				common.ConnWrite("cookie wrong", conn, 20)
				conn.Close()
				continue
			}

			err = common.ConnWrite("ok", conn, 100)
			if err != nil {
				recvLog(ip, err.Error())
				conn.Close()
				continue
			}
			recvLog(ip, "receive connected")
		}  // if conn == nil
		time.Sleep(time.Hour)
		// TODO
	}
}


