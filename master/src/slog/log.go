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
)

var (
	logDir string
	LogChan chan string
)

func Init() {
	logDir = path.Join(config.ConfMap["data-dir"], "log/")

	if !common.IsDirExist(logDir) {
		err := os.MkdirAll(logDir, 0700)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"\nError:\nCan not create log dir %v:\n")
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(3)
		}
	}

	// TODO
}

