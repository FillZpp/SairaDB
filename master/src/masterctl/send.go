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
	"sync/atomic"
	"encoding/json"

	"slog"
	"common"
)

func sendLog(ip, reason string) {
	slog.LogChan<-
		fmt.Sprintf("master controller send task (%v): %v", ip, reason)
}

func sendTask(idx int, ip string, ch chan SendMessage) {
	var conn net.Conn
	var err error
	var msg string
	addr := ip + ":" + port
	buf := make([]byte, 1000)
	once := true
	for {
		if conn == nil {
			conn, err = net.Dial("tcp", addr)
			if err != nil {
				if once {
					sendLog(ip, err.Error())
				}
				once = false
				time.Sleep(time.Second)
				continue
			}
			once = true

			err = common.ConnWrite(cookie, conn, 100)
			if err != nil {
				sendLog(ip, err.Error())
				conn.Close()
				continue
			}

			msg, err = common.ConnRead(buf, conn, 100)
			if err != nil {
				sendLog(ip, err.Error())
				conn.Close()
				continue
			}

			if msg != "ok" {
				sendLog(ip, "wrong cookie")
				conn.Close()
				continue
			}
			atomic.AddInt32(&(MasterList[idx].Status), 1)
			sendLog(ip, "send connected")
		}  // if conn == nil
		
		sm := <-ch
		b, err = json.Marshal(sm.Message)
		if err != nil {
			sm.Ch<- err
			continue
		}
		err = common.ConnWrite(b)
	}
}


