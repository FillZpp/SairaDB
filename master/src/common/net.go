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


package common

import (
	"net"
	"time"
)

func ConnRead(buf []byte, conn net.Conn, timeout int) (ret string, err error) {
	var n int
	for {
		if timeout > 0 {
			conn.SetReadDeadline(time.Now().
				Add(time.Duration(timeout) * time.Millisecond))
		}
		n, err = conn.Read(buf)
		if err != nil {
			return
		}
		
		ret += string(buf[:n])
		if n < len(buf) {
			break
		}
	}
	return
}

func ConnWrite(s string, conn net.Conn, timeout int) (err error) {
	if timeout > 0 {
		conn.SetWriteDeadline(time.Now().
			Add(time.Duration(timeout) * time.Millisecond))
	}
	_, err = conn.Write([]byte(s))
	return
}


