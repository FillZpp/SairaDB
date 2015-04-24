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
	"os"
	"fmt"
	"sync"

	"config"
	"slog"
	"query"
)

type RecvRegister struct {
	id string
	ch chan string
}

type Slave struct {
	ip string
	vnodes []uint64
	
	sendChan chan query.Query
	recvChan chan RecvRegister
	sendStatus int32
	recvStatus int32
}

var (
	port string
	cookie string

	Slaves = make(map[string]*Slave)
	mutex sync.Mutex
)


func Init() {
	vnodeInit()
	
	port, _ = config.ConfMap["slave-port"]
	cookie, _ = config.ConfMap["slave-cookie"]
	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error:\nSlave controller can not listen: %v\n", port)
		os.Exit(3)
	}

	go listenTask(listener)
}

func listenTask(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.LogChan<- "slave controller accept error: " +
				err.Error()
			continue
		}
		go slaveHandler(conn)
	}
}
			
	
