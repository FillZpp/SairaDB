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


package ssignal

import (
	"os"
	"fmt"
	"os/signal"

	"common"
	"meta"
	"slog"
)

var sigChan = make(chan os.Signal, 1)

func Init() {
	sigChan = make(chan os.Signal, 1)
	go sigHanderTask()
	
	signal.Notify(sigChan, os.Interrupt, os.Kill)
}

func sigHanderTask() {
	select {
	case sig := <-sigChan: 
		slog.ExitLog = fmt.Sprintf("by signal %v", sig)
	case s := <-common.ExitChan:
		slog.ExitLog = "because " + s
	}

	meta.ToClose<- true
	<-meta.Closed

	slog.ToClose<- true
	<-slog.Closed

	os.Exit(0)
}

