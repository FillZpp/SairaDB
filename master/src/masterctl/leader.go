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


package masterctl

import (
	"sync/atomic"
	"math/rand"
	"time"
	"fmt"

	"stime"
	"common"
)

func findLeader() {
	for {
		if atomic.LoadInt32(&leader) >= 0 {
			return
		}

		// If already vote for other master
		if !atomic.CompareAndSwapInt32(&voteFor, -1, 0) {
			return
		}
		atomic.AddUint64(&term, 1)
		
		if len(MasterList) == 1 {
			atomic.StoreInt32(&leader, 0)
		}

		allVote := 1
		getVote := 1
		resChans := make([]chan []string, 0, len(MasterList) - 1)
		
		id := stime.GetID()
		v := []string{id, "", "vote", fmt.Sprintf("%v", term)}
		for _, m := range MasterList[1:] {
			// If connection failed
			if atomic.LoadInt32(&(m.Status)) < 2 { 
				continue
			}

			errCh := make(chan error)
			resCh := make(chan []string)
			m.SendChan<- SendMessage{ v, errCh }
			err := <-errCh
			if err != nil {
				continue
			}
			m.RecvChan<- RecvRegister{ id, resCh }
			resChans = append(resChans, resCh)
			allVote++
		}

		for _, resCh := range resChans {
			tmCh := make(chan bool)
			go common.SetTimeout(tmCh, 5000)
			select {
			case res := <-resCh:
				if res[0] == "ok" {
					getVote++
				}
			case <-tmCh:
			}
		}

		atomic.StoreInt32(&voteFor, -1)
		if getVote > allVote/2 {
			fmt.Println("Leader!")
			atomic.StoreInt32(&leader, 0)
			// Tell other masters
			msg := []string{"", "", "leader", fmt.Sprintf("%v", term)}
			for _, m := range MasterList[1:] {
				if atomic.LoadInt32(&(m.Status)) < 2 {
					continue
				}
				m.SendChan<- SendMessage{ msg, make(chan error) }
			}
		}
		
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100) + 100))
	}
}

