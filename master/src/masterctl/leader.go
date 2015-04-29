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


package masterctl

import (
	"sync/atomic"
	"math/rand"
	"fmt"
	"time"

	"stime"
	"meta"
)

func findLeader() {
	for {
		if atomic.LoadInt32(&Leader) >= 0 {
			return
		}

		// If already vote for other master
		if !atomic.CompareAndSwapInt32(&voteFor, -1, 0) {
			return
		}
		meta.Term++
		meta.TermChan<- true
		
		if len(MasterList) == 1 {
			atomic.StoreInt32(&Leader, 0)
		}

		allVote := 1
		getVote := 1
		resChans := make([]chan []string, 0, len(MasterList) - 1)
		
		id := stime.GetID()
		v := []string{id, "vote", fmt.Sprintf("%v", meta.Term)}
		for _, m := range MasterList[1:] {
			// If connection failed
			if atomic.LoadInt32(&(m.Status)) < 2 { 
				continue
			}

			errCh := make(chan error)
			resCh := make(chan []string, 1)
			m.SendChan<- SendMessage{ v, false, errCh }
			err := <-errCh
			if err != nil {
				continue
			}
			m.RecvChan<- RecvRegister{ id, resCh }
			resChans = append(resChans, resCh)
			allVote++
		}

		for _, resCh := range resChans {
			select {
			case res := <-resCh:
				if res[0] == "ok" {
					getVote++
				}
			case <-time.After(5 * time.Second):
			}
		}

		
		if getVote > allVote/2 {
			fmt.Println("Leader!")
			atomic.StoreInt32(&Leader, 0)
			atomic.StoreInt32(&voteFor, -1)
			
			// Tell other masters
			resChans = resChans[:0]
			id = stime.GetID()
			msg := []string{id, "leader", fmt.Sprintf("%v", meta.Term)}
			for _, m := range MasterList[1:] {
				if atomic.LoadInt32(&(m.Status)) < 2 {
					continue
				}

				errCh := make(chan error)
				resCh := make(chan []string, 1)
				m.SendChan<- SendMessage{ msg, false, errCh }
				err := <-errCh
				if err != nil {
					continue
				}
				m.RecvChan<- RecvRegister{ id, resCh }
				resChans = append(resChans, resCh)
			}

			//var otherLeader string
			//var itsFollowerNum int
			for _, resCh := range resChans {
				select {
				case res := <-resCh:
					if res[0] == "ok" {
						atomic.AddInt32(&followerNum, 1)
					} else if res[0] == "no" {
						break
					} else {
						
					}
				case <-time.After(5 * time.Second):
				}
			}
			return
		}

		atomic.StoreInt32(&voteFor, -1)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100) + 100))
	}
}

