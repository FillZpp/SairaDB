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


package clientctl

import (
	"net"
	"encoding/json"
	"fmt"
	"sync/atomic"
	
	"slog"
	"common"
	"query"
	"meta"
)

func handlerLog(ip, reason string) {
	slog.LogChan<-
		fmt.Sprintf("client controller handle (%v): %v", ip, reason)
}

func clientHandler(conn net.Conn) {
	defer conn.Close()
	var msg string
	buf := make([]byte, 1000)
	ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		slog.LogChan<- "client controller address split error: " +
			err.Error()
		return
	}

	msg, err = common.ConnRead(buf, conn, 500)
	if err != nil {
		handlerLog(ip, err.Error())
		return
	}

	if msg != cookie {
		b, _ := json.Marshal([]string{"wrong"})
		handlerLog(ip, "wrong cookie")
		common.ConnWrite(b, conn, 500)
		return
	}

	b, _ := json.Marshal([]string{"ok"})
	err = common.ConnWrite(b, conn, 500)
	if err != nil {
		handlerLog(ip, err.Error())
		return
	}
	
	currentDB := "default"
	for {
		msg, err = common.ConnRead(buf, conn, -1)
		if err != nil {
			handlerLog(ip, err.Error())
			return
		}

		var qry query.CliQuery
		err = json.Unmarshal([]byte(msg), &qry)
		if err != nil {
			fmt.Println("json parse error");
			common.ConnWriteString("(error) parse error", conn, 500)
			continue
		}

		switch qry.Operation {
		case "show_dbs":
			databases := (*map[string]int)(
				atomic.LoadPointer(&(meta.Databases)))
			keys := make([]string, 0, len(*databases))
			for k := range *databases {
				keys = append(keys, k)
			}
			b, _ = json.Marshal(keys)
			common.ConnWrite(b, conn, 500)
		case "create":
			errChan := make(chan error)
			meta.DBChan<- meta.AlterDB{
				errChan,
				"create",
				[]string{qry.Name},
			}
			ret := "ok"
			err = <-errChan
			if err != nil {
				ret = "(error) " + err.Error()
			}
			common.ConnWriteString(ret, conn, 500)
		case "drop":
			errChan := make(chan error)
			meta.DBChan<- meta.AlterDB{
				errChan,
				"drop",
				[]string{qry.Name},
			}
			ret := "ok"
			err = <-errChan
			if err != nil {
				ret = "(error) " + err.Error()
			}
			common.ConnWriteString(ret, conn, 500)
		case "use":
			databases := (*map[string]int)(
				atomic.LoadPointer(&(meta.Databases)))
			ret := "ok"
			if _, ok := (*databases)[qry.Name]; ok {
				currentDB = qry.Name
				fmt.Println("use", currentDB)
			} else {
				ret = fmt.Sprintf("(error) database '%v' does not exist.",
					qry.Name);
			}
			common.ConnWriteString(ret, conn, 500)
			
			// TODO
		default:
			common.ConnWriteString("(error) unknown command", conn, 500)
		}
	}
}

