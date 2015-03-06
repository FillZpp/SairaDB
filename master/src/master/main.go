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


package main

import (
	"fmt"
	"flag"
	"config"
	"os"
	"net"
)

func main() {
	handleFlag()
	
	config.Init()

	// test
	for k, v := range config.ConfMap {
		fmt.Println(k, v)
	}
	fmt.Println(config.MasterList)

	// TODO signal handle
	
}

func handleFlag() {
	help1   := flag.Bool("h", false, "")
	help2   := flag.Bool("help", false, "")
	localIP := flag.String("ip", "", "")
	etcDir := flag.String("etc-dir", "", "")

	flag.Usage = usage
	flag.Parse()

	if *help1 || *help2 {
		flag.Usage()
		os.Exit(0)
	}

	if len(*localIP) > 0 {
		ips, err := net.LookupHost(*localIP)
		if err != nil || len(ips) == 0 || len(ips) > 1 {
			fmt.Fprintf(os.Stderr, "Error:\nsomething wrong with the appointed ip %v", *localIP)
			os.Exit(1)
		}
		config.MasterList = append(config.MasterList, ips[0])
	}

	if len(*etcDir) > 0 {
		config.ConfMap["etc_dir"] = *etcDir
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: saira-master [OPTIONS] flag\n")
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "    -h --help      Display usage")
	fmt.Fprintln(os.Stderr, "    --etc-dir DIR  Find config files in <dir>")
	fmt.Fprintln(os.Stderr, "    --ip IP        Use appoint IP")
}

