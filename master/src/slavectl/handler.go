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

	"slog"
	"common"
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
	buf := make([]byte, 1000)
	ip, _, err = net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		slog.LogChan<- "slave controller address split error: " + err.Error()
		return
	}

	msg, err = common.ConnRead(buf, conn, 100)
	if err != nil {
		handlerLog(ip, err.Error())
		return
	}
	
	if msg != cookie {
		handlerLog(ip, "wrong cookie")
		common.ConnWrite("wrong cookie", conn, 100)
		return
	}
		
	err = common.ConnWrite("ok", conn, 100)
	if err != nil {
		handlerLog(ip, err.Error())
		return
	}
	handlerLog(ip, "connected")

	time.Sleep(time.Hour)
}


