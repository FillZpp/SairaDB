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
use std::sync::atomic::{AtomicUsize, ATOMIC_USIZE_INIT};
use std::fs::{self, PathExt};
use std::io::{stderr, Write};
use std::path::Path;
use super::sr_db::Database;
use super::super::libc;


#[allow(dead_code)]
pub struct VNode {
    id: u64,
    term: u64,
    dups: Vec<String>,
    master: AtomicUsize,
    dbs: HashMap<String, Database>
}

#[allow(dead_code)]
impl VNode {
    pub fn new(id: u64, term: u64) -> VNode {
        VNode {
            id: id,
            term: term,
            dups: Vec::new(),
            master: ATOMIC_USIZE_INIT,
            dbs: HashMap::new()
        }
    }
}

pub fn init(dir: String) -> Vec<u64> {
    let data_dir = dir + "/slave/data/";
    let _ = fs::create_dir_all(Path::new(&data_dir));
    let mut vnodes = Vec::new();

    for entry in match fs::read_dir(data_dir) {
        Ok(a) => a,
        Err(e) => {
            let _ = writeln!(&mut stderr(), "Error:\n read data_dir error:\n{}", e);
            unsafe { libc::exit(3); }
        }
    } {
        match entry {
            Ok(entry) => {
                let path = entry.path();
                if path.is_dir() {
                    if let Some(name) = path.file_name() {
                        if let Some(name) = name.to_str() {
                            let n: Result<u64, _> = name.parse();
                            match n {
                                Ok(n) => vnodes.push(n),
                                Err(_) => {}
                            }
                        }
                    }

                    // TODO
                    // Read vnode term
                }
            }
            Err(e) => {}
        }
    }
    
    vnodes
}

