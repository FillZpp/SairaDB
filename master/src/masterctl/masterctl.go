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
	MasterID int32
	MasterIP string
	SendChan chan SendMessage
	RecvChan chan RecvRegister
	Status int32
}

type SendMessage struct {
	Message []string
	Ever bool
	Ch chan error
}

type RecvRegister struct {
	ID string
	Ret chan []string
}

var (
	port string
	cookie string
	MasterMap = make(map[string]*MasterCtl)
	MasterList []*MasterCtl

	Leader int32 = -1
	voteFor int32 = -1
	followerNum int32 = 0
)

func Init() {
	port, _ = config.ConfMap["master-port"]
	cookie, _ = config.ConfMap["master-cookie"]
	MasterList = make([]*MasterCtl, 1, len(config.MasterList) - 1)
	
	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"\nError:\nMaster controller can not listen :%v\n", port)
		os.Exit(3)
	}

	listenChans := make(map[string]chan net.Conn)
	for idx, ip := range config.MasterList[1:] {
		mc := MasterCtl{
			int32(idx + 1),
			ip,
			make(chan SendMessage, 100),
			make(chan RecvRegister, 100),
			2,
		}
		MasterMap[ip] = &mc
		MasterList = append(MasterList, &mc)
		ch := make(chan net.Conn, 1)
		listenChans[ip] = ch
		go sendTask(idx + 1, ip, mc.SendChan)
		go receiveTask(idx + 1, ip, ch, mc.RecvChan)
	}

	go listenTask(listener, listenChans)
	go findLeader()
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


