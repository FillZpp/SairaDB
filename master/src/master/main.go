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
	"os"
	"path"
	"sync/atomic"
	"crypto/md5"
	"time"

	"common"
	"config"
	"stime"
	"meta"
	"slog"
	"ssignal"

	"csthash"
	"masterctl"
	"slavectl"
)

func main() {
	config.GetHomePath()
	config.Init(getFlags())

	stime.Init()
	meta.Init()
	slog.Init()
	ssignal.Init()

	masterctl.Init()
	slavectl.Init()
	csthash.Init()

	test()
}

func test() {
	for k, v := range config.ConfMap {
		fmt.Println(k, v)
	}
	fmt.Println(config.MasterList)
	fmt.Printf("\n%v\n", csthash.VNodes)
	//meta.DBChan<- meta.AlterDB{ make(chan error, 1), "add_db", []string{"test"}}
	//meta.UserChan<- meta.AlterUser{ make(chan error, 1), "add_user", []string{"wsy", "1234"}}
	time.Sleep(time.Millisecond * 100)
	fmt.Println("term:", atomic.LoadUint64(&(meta.Term)))
	fmt.Println((*map[string]meta.Database)(atomic.LoadPointer(&(meta.Databases))))
	//fmt.Println(*(*string)(atomic.LoadPointer(&(meta.DBEncode))));
	fmt.Println((*map[string]meta.User)(atomic.LoadPointer(&(meta.Users))))
	//fmt.Println(*(*string)(atomic.LoadPointer(&(meta.UserEncode))));

	time.Sleep(time.Millisecond * 100)
}

func getFlags() (flagMap map[string]string) {
	flagMap = make(map[string]string)
	help1   := flag.Bool("h", false, "")
	help2   := flag.Bool("help", false, "")
	
	confDir  := flag.String("conf-dir",
		path.Join(config.Prefix, "/etc/sairadb"), "")
	isLocal  := flag.Bool("local", false, "")
	logLevel := flag.String("log-level", "", "")
	dataDir  := flag.String("data-dir", "", "")
	
	cookie := flag.String("cookie", "", "")

	flag.Usage = usage
	flag.Parse()

	if *help1 || *help2 {
		flag.Usage()
		os.Exit(0)
	}

	if len(*confDir) > 0 {
		if !common.IsDirExist(*confDir) {
			fmt.Fprintf(os.Stderr,
				"Error:\nConfig directory does not exist: %v\n",
				*confDir)
			os.Exit(2)
		}
		flagMap["conf-dir"] = *confDir
	}
	
	if *isLocal {
		flagMap["local"] = "on"
	}

	if len(*logLevel) > 0 {
		switch *logLevel {
		case "error":  fallthrough
		case "slow":   fallthrough
		case "full":   flagMap["log-level"] = *logLevel
		default: {
			fmt.Fprintf(os.Stderr,
				"Error:\nInvalid given log-level: %v\n",
				*logLevel)
			os.Exit(2)
		}
		}
	}

	if len(*dataDir) > 0 {
		flagMap["data-dir"] = *dataDir
	}

	if len(*cookie) > 0 {
		ckMd5 := fmt.Sprintf("%x", md5.Sum([]byte(*cookie)))
		flagMap["client-cookie"] = ckMd5
		flagMap["master-cookie"] = ckMd5
		flagMap["slave-cookie"] = ckMd5
	}

	return flagMap
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: saira-master [OPTIONS]\n")
	fmt.Fprintln(os.Stderr, "Options:")
	
	fmt.Fprintln(os.Stderr,
		"    --conf-dir DIR      Find config files in <DIR>")
	fmt.Fprintln(os.Stderr,
		"    --local             Cluster only on local machine")
	fmt.Fprintln(os.Stderr,
		"    --log-level LEVEL   Define log level [error/slow/full]")
	fmt.Fprintln(os.Stderr,
		"    --data-dir DIR      Save meta data and log in <DIR>/master")
	fmt.Fprintln(os.Stderr,
		"    --cookie COOKIE     Set cookie for connection safety")
	fmt.Fprintln(os.Stderr,
		"    -h --help           Print this help menu")
}


