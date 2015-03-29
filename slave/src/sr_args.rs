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


extern crate libc;
extern crate getopts;

use std::io::{stderr, Write};
use std::collections::HashMap;
use std::path::Path;
use std::fs::PathExt;
use std::env;
use self::getopts::Options;


pub fn get_flags() -> HashMap<String, String> {
    let mut flags = HashMap::new();
    let args: Vec<String> = env::args().collect();
    
    let mut opts = Options::new();
    opts.optflag("h", "help", "Print this help menu");
    opts.optopt("", "conf-dir", "Find config files in <DIR>", "DIR");

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

    match matches.opt_str("conf-dir") {
        Some(a) => {
            if !Path::new(&a).is_dir() {
                let _ = writeln!(&mut stderr(),
                                 "Error:\nConfig directory does not exist: {}",
                                 a);
                unsafe { libc::exit(3); }
            }
            flags.insert("conf-dir".to_string(), a);
        }
        None => {}
    }

    flags
}

fn usage(opts: Options) {
    print!("{}", opts.usage("Usage: saira-slave [OPTIONS]"));
}

