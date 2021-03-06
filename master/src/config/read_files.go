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


package config

import (
	"os"
	"path"
	"io/ioutil"
	"net"
	"fmt"
	"strconv"
	"crypto/md5"
)

var (
	localIPs []string
	mastersFile string
)

func readMastersFile() {
	host, err := os.Hostname()
	if err != nil {
		fmt.Fprintln(os.Stderr,
			"Error: can not get local hostname.")
		os.Exit(3)
	}
	
	localIPs, err = net.LookupHost(host)
	if err != nil || len(localIPs) == 0 {
		fmt.Fprintln(os.Stderr, "Error: can not find local ip.")
		os.Exit(2)
	}

	mastersFile = path.Join(ConfMap["conf-dir"] + "/masters")
	fl, err := os.Open(mastersFile)
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
					insertIP(string(con[prev:i]), lineNum)
					step = 2
				}
				continue
			} else if con[i] == '\n' {
				if step == 1 {
					insertIP(string(con[prev:i]), lineNum)
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
			if step == 1 {
				insertIP(string(con[prev:i]), lineNum)
			}
			break
		}
		i++
	}

	if MasterList[0] == "" {
		fmt.Fprintf(os.Stderr,
			"Error:\nNo local ip appointed in %v\n",
			mastersFile)
		os.Exit(3)
	}
}

func insertIP(host string, lineNum int) {
	ips, err := net.LookupHost(host)
	if err != nil || len(ips) == 0 {
		fmt.Fprintf(os.Stderr,
			"Error:\n%v line %d: %v\nCan not find ip from the hostname\n",
			mastersFile, lineNum, host)
		os.Exit(3)
	} else if len(ips) > 1 {
		fmt.Fprintln(os.Stderr,
			"Error:\n%v line %d: %v\nFind more than one ip from the hostname.\n",
			mastersFile, lineNum, host)
		os.Exit(3)
	}
	
	for _, v := range localIPs {
		if ips[0] == v {
			if MasterList[0] != "" {
				fmt.Fprintf(os.Stderr,
					"Error:\nMore than one local ips in file %v\n",
					mastersFile)
				os.Exit(3)
			}
			MasterList[0] = v
			return
		}
	}

	for _, v := range MasterList {
		if ips[0] == v {
			return
		}
	}
			
	MasterList = append(MasterList, ips[0])
}

func readConfFile(confDir string) {
	fl, err := os.Open(path.Join(confDir + "/master.conf"))
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
			if step == 3 {
				updateConf(key, string(con[prev:i]), lineNum)
			}
			break
		}
		i++
	}
}

func updateConf(key, value string, lineNum int) {
	_, ok := ConfMap[key]
	if !ok {
		fmt.Fprintf(os.Stderr,
			"Error:\n%v line %d: %v %v\nUnknown config.\n",
			path.Join(ConfMap["conf-dir"] + "/master.conf"), lineNum,
			key, value)
		os.Exit(3)
	}

	for _, v := range BoolConfs {
		if key == v {
			if !(value == "on" || value == "off") {
				fmt.Fprintln(os.Stderr,
					"Error:\nInvalid config, this must be 'on' or 'off':\n")
				fmt.Fprintf(os.Stderr, "%v %v\n", key, value)
				os.Exit(3)
			}
			ConfMap[key] = value
			return
		}
	}

	if key == "log-level" {
		switch value {
		case "error":  fallthrough
		case "slow":   fallthrough
		case "full":   break
		default: {
			fmt.Fprintf(os.Stderr,
				"Error:\nInvalid config:\nlog-level %v\n",
				value)
			os.Exit(3)
		}
		}
		ConfMap[key] = value
		return
	}

	for _, v := range UintConfs {
		if key == v {
			i, err := strconv.Atoi(value)
			if err != nil || i < 1 {
				fmt.Fprintf(os.Stderr,
					"Error:\nInvalid config:\n%v %v\n",
					key, value)
				os.Exit(3)
			}
			ConfMap[key] = value
			return
		}
	}

	if key == "client-cookie" || key == "master-cookie" ||
		key == "slave-cookie" {
		ConfMap[key] = fmt.Sprintf("%x", md5.Sum([]byte(value)))
		return
	}

}

