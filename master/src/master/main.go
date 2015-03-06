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
	"path"
)

func main() {
	config.GetHomePath()
	flagMap := handleFlag()
	
	config.Init(flagMap)

	// test
	for k, v := range config.ConfMap {
		fmt.Println(k, v)
	}
	fmt.Println(config.MasterList)

	// TODO signal handle
	
}

func handleFlag() (flagMap map[string]string) {
	flagMap = make(map[string]string)
	help1   := flag.Bool("h", false, "")
	help2   := flag.Bool("help", false, "")
	
	confDir   := flag.String("conf-dir",
		path.Join(config.Prefix, "/etc/sairadb"), "")
	ip       := flag.String("ip", "", "")
	isLocal  := flag.Bool("local", false, "")
	logLevel := flag.String("log-level", "common", "")
	dataDir  := flag.String("data-dir",
		path.Join(config.HomeDir, "/saira_data"), "")

	flag.Usage = usage
	flag.Parse()

	if *help1 || *help2 {
		flag.Usage()
		os.Exit(0)
	}

	if len(*confDir) > 0 {
		if !isDirExist(*confDir) {
			fmt.Fprintf(os.Stderr,
				"\nError:\nInvalid directory: %v\n",
				*confDir)
			os.Exit(2)
		}
		flagMap["conf-dir"] = *confDir
	}
	
	if *isLocal {
		flagMap["local"] = "on"
	} else {
		if len(*ip) > 0 {
			ips, err := net.LookupHost(*ip)
			if err != nil || len(ips) == 0 || len(ips) > 1 {
				fmt.Fprintf(os.Stderr,
					"\nError:\nSomething wrong with the appointed ip: %v\n",
					*ip)
				os.Exit(2)
			}
			config.MasterList = append(config.MasterList, ips[0])
		}
	}

	switch *logLevel {
	case "error":  fallthrough
	case "common": fallthrough
	case "slow":   fallthrough
	case "full":   flagMap["log-level"] = *logLevel
	default: {
		fmt.Fprintf(os.Stderr,
			"\nError:\nInvalid given log-level: %v\n", *logLevel)
		os.Exit(2)
	}
	}

	flagMap["data-dir"] = *dataDir
	
	return flagMap
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: saira-master [OPTIONS] flag\n")
	fmt.Fprintln(os.Stderr, "Options:")
	
	fmt.Fprintln(os.Stderr,
		"    --conf-dir DIR     Find config files in <dir>")
	fmt.Fprintln(os.Stderr,
		"    --local            Clusters only on local machine")
	fmt.Fprintln(os.Stderr,
		"    --ip IP            Use appoint IP")
	fmt.Fprintln(os.Stderr,
		"    --log-level LEVEL  Define log level [error/common/slow/full]")
	fmt.Fprintln(os.Stderr,
		"    --data-dir DIR     Save meta data and log in {DIR}/master")
	fmt.Fprintln(os.Stderr,
		"    -h --help          Display usage")
}

func isDirExist(dir string) bool {
	fl, err := os.Stat(dir)
	if err != nil {
		return os.IsExist(err)
	}
	return fl.IsDir()
}


