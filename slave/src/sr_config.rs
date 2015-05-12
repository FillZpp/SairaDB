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


use std::io::{stderr, Write, Read};
use std::collections::HashMap;
use std::fs::File;
use std::net;
use super::libc;


pub fn init(map: &HashMap<String, String>)  -> Vec<String> {
    read_masters(map.get("conf-dir").unwrap().to_string())
}

fn read_masters(conf_dir: String) -> Vec<String> {
    let mut masters = Vec::new();

    let fname = conf_dir + "/masters";
    let mut f = match File::open(&fname) {
        Ok(f) => f,
        Err(e) => {
            let _ = writeln!(&mut stderr(),
                             "Error:\nCan not find {}:\n{}", fname, e);
            unsafe { libc::exit(3); }
        }
    };

    let mut s = String::new();
    if let Err(e) = f.read_to_string(&mut s) {
        let _ = writeln!(&mut stderr(), "Error:\nRead {}:\n{}", fname, e);
        unsafe { libc::exit(3); }
    }

    let lines = s.split('\n');
    let mut n = 0;
    for line in lines {
        n += 1;
        let host = line.trim();
        if host.len() > 0 {
            let mut ips = match net::lookup_host(host) {
                Ok(a) => a,
                Err(e) => {
                    let _ = writeln!(&mut stderr(),
                                     "Error:\n{} line {}: {}\n{}",
                                     fname, n, host, e);
                    unsafe { libc::exit(3); }
                }
            };

            masters.push(match ips.next() {
                Some(Ok(a)) => format!("{}", a.ip()),
                Some(Err(e)) => {
                    let _ = writeln!(&mut stderr(),
                                     "Error:\n{} line {}: {}\n{}",
                                     fname, n, host, e);
                    unsafe { libc::exit(3); }
                }
                None => {
                    let _ = writeln!(&mut stderr(),
                                     "Error:\n{} line {}: {}\nCan not find ip",
                                     fname, n, host);
                    unsafe { libc::exit(3); }
                }
            });
        }
    }

    masters
}


