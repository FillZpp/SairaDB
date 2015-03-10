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


// Package config generate CFGMap from configure file and default configures.
package config

import (
	"strings"
	"os"
	"fmt"
)

var (
	ConfMap = make(map[string]string)
	MasterList = make([]string, 0, 5)
	HomeDir string
)

// Initialize configure for only once
func Init(flagMap map[string]string) {
	initConf()

	readConfFile(flagMap["conf-dir"])
	for k, v := range flagMap {
		ConfMap[k] = v
	}
	
	local, _ := ConfMap["local"]
	if local != "on" {
		readMastersFile()
	}
}

func GetHomePath() {
	for _, env := range os.Environ() {
		if strings.Index(env, "HOME=") == 0 ||
			strings.Index(env, "HOMEPATH=") == 0 {
			HomeDir = strings.SplitN(env, "=", 2)[1]
		}
	}
	
	if HomeDir == "" {
		fmt.Fprintln(os.Stderr, "\nError:\nCan not find HOME path")
		os.Exit(2)
	}
}


