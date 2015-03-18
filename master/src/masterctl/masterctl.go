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
	"os"
	"fmt"
	
	"config"
	"slog"
)

type MasterCtl struct {
	Send chan string
	Receive chan string
}

type AlterMaster struct {
	Ch chan error
	AlterType string
	AlterCont string
}

var (
	MasterMap = make(map[string]MasterCtl)
	AlterChan = make(chan AlterMaster, 10)
)

func Init() {
	listener, err := net.Listen("tcp", ":" + config.ConfMap["master-port"])
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"\nError:\nCan not listen :%v\n",
			config.ConfMap["master-port"])
		os.Exit(3)
	}

	listenChans := make(map[string]chan net.Conn)
	for _, ip := range config.MasterList {
		mc := MasterCtl{ make(chan string), make(chan string) }
		MasterMap[ip] = mc
		ch := make(chan net.Conn, 1)
		listenChans[ip] = ch
		go sendTask(ip, mc.Send)
		go receiveTask(ip, ch)
	}

	go listenTask(listener, listenChans)
}

func listenTask(listener net.Listener, listenChans map[string]chan net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.LogChan<- "master controller accept error: " +
				err.Error()
			continue
		}
		
		ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			slog.LogChan<- "master controller address split error: " +
				err.Error()
			conn.Close()
			continue
		}

		ch, ok := listenChans[ip]
		if !ok {
			slog.LogChan<-
				fmt.Sprintf("master controller get %v is not in masters",
				ip)
			conn.Close()
			continue
		}

		ch<- conn
	}
}


