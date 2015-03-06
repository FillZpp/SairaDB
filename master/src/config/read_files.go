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


package config

import (
	"os"
	"path"
	"io/ioutil"
	"net"
	"fmt"
)

func readMastersFile() {
	if len(MasterList) == 0 {
		host, err := os.Hostname()
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"\nError: can not get local hostname")
			os.Exit(3)
		}
		
		localIPs, err := net.LookupHost(host)
		if err != nil || len(localIPs) == 0 {
			fmt.Fprintln(os.Stderr, "\nError: can not find local ip.")
			os.Exit(2)
		} else if len(localIPs) > 1 {
			fmt.Fprintln(os.Stderr,
				"\nError: Find more than one ip for local hostname.\n",
				"You should use option [-ip] to appoint one.")
			os.Exit(3)
		}
		MasterList = append(MasterList, localIPs[0])
	}

	fl, err := os.Open(path.Join(ConfMap["etc_dir"] + "masters"))
	if err != nil {
		return
	}
	defer fl.Close()

	con, err := ioutil.ReadAll(fl)
	if err != nil {
		return 
	}

	len := len(con)
	i := 0
	lineNum := 0
	for {
		lineNum++
		prev := i
		step := 0
		for ; i < len; i++ {
			if con[i] == ' ' {
				if step == 1 {
					checkIP(string(con[prev:i]), lineNum)
					step = 2
				}
				continue
			} else if con[i] == '\n' {
				if step == 1 {
					checkIP(string(con[prev:i]), lineNum)
				}
				break
			} else {
				if step == 0 {
					if con[i] == '#' {
						step = 2
					} else {
						prev = i
						step = 1
					}
				}
			}
		}
		if i == len {
			break
		}
		i++
	}
}

func checkIP(host string, lineNum int) {
	ips, err := net.LookupHost(host)
	if err != nil || len(ips) == 0 {
		fmt.Fprintln(os.Stderr,
			"\nError:")
		fmt.Fprintf(os.Stderr,
			"%v line %d: %v\nCan not find ip from the hostname\n",
			path.Join(ConfMap["etc_dir"] + "masters"), lineNum, host)
		os.Exit(3)
	} else if len(ips) > 1 {
		fmt.Fprintln(os.Stderr,
			"\nError:")
		fmt.Fprintln(os.Stderr,
			"%v line %d: %v\nFind more than one ip from the hostname.\n",
			path.Join(ConfMap["etc_dir"] + "masters"), lineNum, host)
		os.Exit(3)
	}
	
	for _, m := range MasterList {
		if ips[0] == m {
			return
		}
	}
	MasterList = append(MasterList, ips[0])
}

func readConfFile() {
	fl, err := os.Open(path.Join(ConfMap["etc_dir"] + "master.conf"))
	if err != nil {
		return
	}
	defer fl.Close()

	con, err := ioutil.ReadAll(fl)
	if err != nil {
		return
	}

	len := len(con)
	i := 0
	lineNum := 0
	for {
		lineNum++
		prev := i
		step := 0
		var key string
		for ; i < len; i++ {
			if con[i] == ' ' {
				if step == 1 {
					key = string(con[prev:i])
					step = 2
				} else if step == 3 {
					updateConf(key, string(con[prev:i]), lineNum)
					step = 4
				}
				continue
			} else if con[i] == '\n' {
				if step == 3 {
					updateConf(key, string(con[prev:i]), lineNum)
				}
				break
			} else {
				if step == 0 {
					if con[i] == '#' {
						step = 4
					} else {
						prev = i
						step = 1
					}
				} else if step == 2 {
					prev = i
					step = 3
				}
			}
		}
		if i == len {
			break
		}
		i++
	}
}

func updateConf(key, value string, lineNum int) {
	_, ok := ConfMap[key]
	if !ok {
		fmt.Fprintln(os.Stderr,
			"Error:")
		fmt.Fprintf(os.Stderr,
			"%v line %d: %v %v\nUnknown config.\n",
			path.Join(ConfMap["etc_dir"] + "master.conf"), lineNum,
			key, value)
		os.Exit(3)
	}
	ConfMap[key] = value
}

