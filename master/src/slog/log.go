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


package slog

import (
	"fmt"
	"os"
	"path"

	"common"
	"config"
	"stime"
)

var (
	logDir string
	LogChan = make(chan string)
	
	logFile *os.File
	fileSize int

	ToClose = make(chan bool)
	Closed = make(chan bool)
	ExitLog string
)

func Init() {
	logDir = path.Join(config.ConfMap["data-dir"],
		"master/log/" + stime.TimeCat())

	if !common.IsDirExist(logDir) {
		err := os.MkdirAll(logDir, 0700)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Error:\nCan not create log dir %v:\n%v\n",
				err.Error())
			os.Exit(3)
		}
	}

	newLogFile(true)
	go task()
}

func newLogFile(b bool) {
	fname := path.Join(logDir,
		"/" + stime.TimeCat() + ".log")
	lf, err := os.OpenFile(fname, os.O_RDWR | os.O_CREATE, 0600)
	if err != nil {
		if b {
			fmt.Fprintf(os.Stderr,
				"Error:\nCan not create new log file %v:\n%v\n",
				fname,
				err.Error())
			os.Exit(3)
		} else {
			common.ExitChan<- fmt.Sprintf("Error: Can not create new log file %v: %v",
				fname,
				err.Error())
		}
	}
	logFile.Close()
	logFile = lf
	fileSize = 0
}

func task() {
	cache := ""
	size := 0
	var newLog string
	for {
		if size == 0 {
			select {
			case <-ToClose:
				logFile.WriteString(stime.TimeFormat() + 
					" This master closed " + ExitLog + ".\n")
				logFile.Close()
				Closed<- true
				for {
					<-LogChan
				}
			case newLog = <-LogChan:
			}
			cache = stime.TimeFormat() + " " + newLog + "\n"
			size++
		} else {
			if size < 100 {
				ch := make(chan bool)
				go common.SetTimeout(ch, 100)
				select {
				case newLog = <-LogChan:
					cache = stime.TimeFormat() +
						" " + newLog + "\n"
					size++
					continue
				case <-ch:
				}
			}
			
			l, err := logFile.WriteString(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"\nError:\nWrite log file error:\n%v\n",
					err.Error())
				os.Exit(4)
			}
			
			fileSize += l
			cache = ""
			size = 0

			if fileSize >= 1e9 {
				newLogFile(false)
			}
		}
	}
}

