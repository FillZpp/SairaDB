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
	"time"
	"sync/atomic"
	"encoding/json"
	"container/list"

	"slog"
	"common"
	"meta"
)

func sendLog(ip, reason string) {
	slog.LogChan<-
		fmt.Sprintf("master controller send task (%v): %v", ip, reason)
}

func sendTask(idx int, ip string, sendChan chan SendMessage) {
	var conn net.Conn
	var err error
	var msg string
	var b []byte
	var sm SendMessage
	addr := ip + ":" + port
	buf := make([]byte, 10)
	onceLog := true
	waitSendList := list.New()
	for {
		if conn == nil {
			conn, err = net.Dial("tcp", addr)
			if err != nil {
				if onceLog {
					sendLog(ip, err.Error())
					atomic.AddInt32(&(MasterList[idx].Status), -1)
					a := true
					for a {
						select {
						case sm = <-sendChan:
							sm.Ch<- err
							if sm.Ever {
								waitSendList.PushFront(sm)
							}
						case <-time.After(100 * time.Millisecond):
							a = false
						}
					}
				}
				onceLog = false
				time.Sleep(time.Second * 5)
				continue
			}
			onceLog = true

			err = common.ConnWriteString(cookie, conn, 500)
			if err != nil {
				sendLog(ip, err.Error())
				conn.Close()
				conn = nil
				continue
			}

			msg, err = common.ConnRead(buf, conn, 500)
			if err != nil {
				sendLog(ip, err.Error())
				conn.Close()
				conn = nil
				continue
			}

			if msg != "ok" {
				sendLog(ip, "wrong cookie")
				conn.Close()
				conn = nil
				continue
			}
			atomic.AddInt32(&(MasterList[idx].Status), 1)
			sendLog(ip, "send connected")

			if atomic.LoadInt32(&Leader) == 0 {
				b, _ = json.Marshal([]string{"",
					"leader", fmt.Sprintf("%v",
						atomic.LoadUint64(&(meta.Term)))})
				err = common.ConnWrite(b, conn, 500)
				if err != nil {
					sendLog(ip, err.Error())
					conn.Close()
					conn = nil
					continue
				}

				_, err = common.ConnRead(buf, conn, 500)
				if err != nil {
					sendLog(ip, err.Error())
					conn.Close()
					conn = nil
					continue
				}
			}
		}  // if conn == nil

		ele := waitSendList.Front()
		for {
			if ele != nil {
				sm = ele.Value.(SendMessage)
			} else {
				select {
				case <-time.After(time.Second):
					sm = SendMessage{
						[]string{"heartbeat"},
						false,
						make(chan error, 1),
					}
				case sm = <-sendChan:
				}
			}
			
			b, _ = json.Marshal(sm.Message)
			
			err = common.ConnWrite(b, conn, 500)
			if err == nil {
				_, err = common.ConnRead(buf, conn, 500)
			}
			if err != nil {
				sm.Ch<- err
				sendLog(ip, err.Error())
				conn.Close()
				conn = nil
				if ele == nil && sm.Ever {
					waitSendList.PushFront(sm)
				}
				break
			}
			sm.Ch<- nil

			if ele == nil {
				ele = waitSendList.Front()
			} else {
				waitSendList.Remove(ele)
				ele.Next()
			}
		}
	}
}


