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


extern crate crypto;
extern crate getopts;
extern crate libc;

use std::io::{stderr, Write};
use std::collections::HashMap;
use std::path::Path;
use std::fs::{self, PathExt};
use std::env;
use self::getopts::Options;
use self::crypto::digest::Digest;
use self::crypto::md5::Md5;


pub fn get_flags() -> HashMap<String, String> {
    let mut flag_map = HashMap::new();
    let args: Vec<String> = env::args().collect();
    
    let mut opts = Options::new();
    opts.optflag("h", "help", "Print this help menu");
    opts.optopt("", "conf-dir", "Find config files in <DIR>", "DIR");
    opts.optopt("", "data-dir", "Save data in <DIR>/slave", "DIR");
    opts.optopt("", "cookie", "Set cookie for connection safety", "COOKIE");
    opts.optopt("p", "port", "Listen port", "PORT");
    opts.optopt("", "master-port", "Port master listened for slave", "PORT");

    let matches = match opts.parse(&args[1..]) {
        Ok(m) => m,
        Err(f) => {
            let _ = writeln!(&mut stderr(), "Error:\nArgs parse error:\n{}", f);
            unsafe { libc::exit(3); }
        }
    };

    if matches.opt_present("h") {
        usage(opts);
        unsafe { libc::exit(0); }
    }

    flag_map.insert("conf-dir".to_string(), {
        let dir = match matches.opt_str("conf-dir") {
            Some(a) => a,
            None => super::sr_prefix::PREFIX.to_string() + "/etc/sairadb"
        };
        
        if !Path::new(&dir).is_dir() {
            let _ = writeln!(&mut stderr(),
                             "Error:\nConfig directory does not exist: {}",
                             dir);
            unsafe { libc::exit(3); }
        }
        dir
    });

    flag_map.insert("data-dir".to_string(), {
        let dir = match matches.opt_str("data-dir") {
            Some(a) => a,
            None => match env::home_dir() {
                Some(ref p) => format!("{}/saira_data", p.display()),
                None => {
                    let _ = writeln!(&mut stderr(),
                                     "Error:\nCan not get HOME directory");
                    unsafe { libc::exit(3); }
                }
            }
        };
        let _ = fs::create_dir_all(Path::new(&dir));
        dir
    });

    flag_map.insert("cookie".to_string(), {
        let ck = match matches.opt_str("cookie") {
            Some(a) => a,
            None => "".to_string()
        };
        let mut cal = Md5::new();
        cal.input_str(&ck);
        cal.result_str().to_string()
    });

    flag_map.insert("port".to_string(), match matches.opt_str("port") {
        Some(a) => a,
        None => "4403".to_string()
    });

    flag_map.insert("master-port".to_string(), match matches.opt_str("master-port") {
        Some(a) => a,
        None => "4402".to_string()
    });

    flag_map
}

fn usage(opts: Options) {
    print!("{}", opts.usage("Usage: saira-slave [OPTIONS]"));
}

