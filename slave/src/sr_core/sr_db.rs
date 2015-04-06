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


extern crate serialize;

use std::collections::{BTreeMap, HashMap, VecDeque};
use std::sync::atomic::{AtomicUsize, ATOMIC_USIZE_INIT, Ordering};
use std::sync::mpsc::{channel, Sender, Receiver};
use std::sync::{Arc, Mutex, RwLock};
use std::thread;
use std::ops::Drop;
use self::serialize::json::Json;


#[allow(dead_code)]
pub struct Query {
    ty: String,
}

#[allow(dead_code)]
struct Unit {
    data: Arc<Json>,
    mtx: Mutex<bool>  // lock when change data
}

#[allow(dead_code)]
struct Page {
    id: u64,
    size: AtomicUsize,
    units: Arc<BTreeMap<String, Unit>>,
    // read lock when change Unit in it , write lock when add/delete Unit
    rwlock: RwLock<bool>  
}

#[allow(dead_code)]
struct SetContent {
    key: String,
    size: AtomicUsize,
    page_id: AtomicUsize,
    pages: VecDeque<Arc<Page>>,
}

#[allow(dead_code)]
struct Set {
    name: String,
    sender: Sender<Query>,
}

#[allow(dead_code)]
pub struct Database {
    name: String,
    sender: Sender<Query>,
    sets: HashMap<String, Set>
}

#[allow(dead_code)]
impl Unit {
    pub fn new(j: Json) -> Unit {
        Unit {
            data: Arc::new(j),
            mtx: Mutex::new(true)
        }
    }
}

#[allow(dead_code)]
impl Page {
    pub fn new(id: u64) -> Page {
        Page {
            id: id,
            size: ATOMIC_USIZE_INIT,
            units: Arc::new(BTreeMap::new()),
            rwlock: RwLock::new(true)
        }
    }
}

#[allow(dead_code)]
impl Set {
    pub fn new(name: String, key: String, log_sender: Sender<String>) -> Set {
        let mut vd = VecDeque::new();
        vd.push_front(Arc::new(Page::new(0)));
        let (tx, rx) = channel();
        let set_cont = SetContent {
            key: key,
            size: ATOMIC_USIZE_INIT,
            page_id: AtomicUsize::new(1),
            pages: vd,
        };
        
        thread::spawn(move || {
            set_task(set_cont, rx, log_sender);
        });
    
        Set {
            name: name,
            sender: tx,
        }
    }
}

#[allow(dead_code)]
struct PageThreadStatus {
    id: u32,
    task_num: Arc<AtomicUsize>,
    sender: Sender<Query>
}

#[allow(dead_code)]
fn page_task(set_cont: Arc<SetContent>, receiver: Receiver<Query>,
             task_num: Arc<AtomicUsize>, log_sender: Sender<String>) {
    loop {
        let qr = match receiver.recv() {
            Ok(qr) => qr,
            Err(e) => {
                let _ = log_sender.send(
                    format!("slave core page_task receive error: {}", e));
                continue;
            }
        };
    }
}

#[allow(dead_code)]
fn set_task(set_cont: SetContent, receiver: Receiver<Query>,
            log_sender: Sender<String>) {
    let set_cont = Arc::new(set_cont);
    let td_num = unsafe { super::td_num };
    let mut p_tasks = Vec::with_capacity(td_num as usize);
    for i in 0..td_num {
        let (tx, rx) = channel();
        let set_cont = set_cont.clone();
        let task_num = Arc::new(ATOMIC_USIZE_INIT);
        let task_num_clone = task_num.clone();
        let log_sender = log_sender.clone();
        thread::spawn(move || {
            page_task(set_cont, rx, task_num_clone, log_sender);
        });

        p_tasks.push(PageThreadStatus {
            id: i,
            task_num: task_num,
            sender: tx
        });
    }
    let cur_thread = 0;

    loop {
        let qr = match receiver.recv() {
            Ok(qr) => qr,
            Err(e) => {
                let _ = log_sender.send(
                    format!("slave core set_task receive error: {}", e));
                continue;
            }
        };

        // TODO
    }
}

#[allow(dead_code)]
fn find_least_load(p_tasks: &Vec<PageThreadStatus>) -> usize {
    let mut idx = -1;
    let mut min = -1;
    let len = p_tasks.len();
    
    for i in 0..len {
        let n = p_tasks[i].task_num.load(Ordering::Relaxed);
        if n == 0 {
            idx = i;
            break;
        } else if n < min {
            idx = i;
            min = n;
        }
    }
    return idx;
}


