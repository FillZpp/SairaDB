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


use std::collections::HashMap;
use std::net::TcpStream;
use std::thread;
use std::sync::Arc;
use std::sync::mpsc::Sender;
use std::sync::atomic::AtomicUsize;
use super::libc;


pub fn init(masters: Vec<String>, vnodes: Vec<u64>,
            map: &HashMap<String, String>, log_sender: Sender<String>) {
    let vnodes = Arc::new(vnodes);
    let cookie = map.get("cookie-master").unwrap().to_string();
    let port = map.get("master-port").unwrap().to_string();
    for master in masters {
        let vnodes = vnodes.clone();
        let log_sender = log_sender.clone();
        let master = master.to_string();
        let cookie = cookie.to_string();
        let port = port.to_string();
        let _ = thread::Builder::new().name(format!("master_task({})", master))
            .spawn(|| {
                master_task(master, port, vnodes, cookie, log_sender);
        });
    }
}

fn master_task(ip: String, port: String, vnodes: Arc<Vec<u64>>, cookie: String,
               log_sender: Sender<String>) {
    let addr: &str = &(ip + ":" + &port);
    let count = Arc::new(AtomicUsize::new(0));
    loop {
        
        let stream = TcpStream::connect(addr);
        //match stream.write_all(cookie.as_bytes());
    }
}

fn master_connection() {
}

