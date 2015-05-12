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


extern crate time;
extern crate libc;

use std::thread;
use std::sync::mpsc::{channel, Sender, Receiver};
use std::io::{stderr, Write};
use std::fs::{create_dir_all, File};
use std::path::Path;


pub fn init(dir: String) -> (Sender<String>, thread::JoinHandle<()>) {
    let tm = time::now();
    let log_dir = dir + "/slave/log/" + &tm_cat(tm);
    let _ = create_dir_all(Path::new(&log_dir));

    let (tx, rx) = channel();
    let jh = thread::Builder::new().name("log".to_string()).spawn(|| {
        log_task(log_dir, rx);
    }).unwrap();
    (tx, jh)
}

fn tm_cat(tm: time::Tm) -> String {
    format!("{}-{}-{}_{}-{}-{}",
            1900 + tm.tm_year, 1 + tm.tm_mon, tm.tm_mday,
            tm.tm_hour, tm.tm_min, tm.tm_sec)
}

fn tm_format(tm: time::Tm) -> String {
    format!("{}-{}-{} {}:{}:{}",
            1900 + tm.tm_year, 1 + tm.tm_mon, tm.tm_mday,
            tm.tm_hour, tm.tm_min, tm.tm_sec)
}

fn log_task(log_dir: String, rx: Receiver<String>) {
    let mut tm = time::now();
    let mut fname = log_dir.to_string() + "/" + &tm_cat(tm) + ".log";
    let mut log_file = match File::create(Path::new(&fname)) {
        Ok(f) => f,
        Err(e) => {
            let _ = writeln!(&mut stderr(),
                             "Error:\nCan not create {}:\n{}", fname, e);
            unsafe { libc::exit(4); }
        }
    };
    let mut size = 0usize;
    let mut cache = "".to_string();
    //let mut timer = Timer::new().unwrap();
    
    loop {
        if cache.len() == 0 {
            let new_log = rx.recv().unwrap();
            tm = time::now();
            cache = tm_format(tm) + " " + &new_log + "\n";
        } else {
            //let timeout = timer.oneshot(Duration::milliseconds(100));
            /*select! {
                new_log = rx.recv() => {
                    cache = cache + &tm_format(tm) + " " + &new_log.unwrap() + "\n";
                    continue;
                }
                _ = timeout.recv() => {}
            }*/

            match log_file.write_all(cache.as_bytes()) {
                Ok(_) => {println!("ok");}
                Err(e) => {
                    let _ = writeln!(&mut stderr(),
                                     "Error:\nWrite {}:\n{}", fname, e);
                    unsafe { libc::exit(4); }
                }
            }
            size += cache.len();
            cache = "".to_string();

            if size > 1000000000 {
                tm = time::now();
                fname = log_dir.to_string() + "/" + &tm_cat(tm) + ".log";
                log_file = match File::create(Path::new(&fname)) {
                    Ok(f) => f,
                    Err(e) => {
                        let _ = writeln!(&mut stderr(),
                                         "Error:\nCan not create {}:\n{}",
                                         fname, e);
                        unsafe { libc::exit(4); }
                    }
                };
                size = 0;
            }
        }
    }
}
    
