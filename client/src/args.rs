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


extern crate crypto;
extern crate getopts;
extern crate libc;

use std::io::{stderr, Write};
use std::collections::HashMap;
use std::env;
use self::getopts::Options;
use self::crypto::digest::Digest;
use self::crypto::md5::Md5;


pub fn get_flags() -> HashMap<String, String> {
    let mut flag_map = HashMap::new();
    let args: Vec<String> = env::args().collect();

    let mut opts = Options::new();
    opts.optflag("v", "version", "Print version");
    opts.optflag("h", "help", "Print this help menu");
    opts.optopt("", "cookie", "Set cookie for connection safety", "COOKIE");
    opts.optopt("", "ip", "Master IP", "IP");
    opts.optopt("", "master-port", "Port master listened for client", "PORT");
    opts.optopt("", "slave-port", "Port slave listened for client", "PORT");

    let matches = match opts.parse(&args[1..]) {
        Ok(m) => m,
        Err(f) => {
            let _ = writeln!(&mut stderr(), "Error:\nArgs parse error:\n{}", f);
            unsafe { libc::exit(3); }
        }
    };

    if matches.opt_present("version") {
        print_version();
        unsafe { libc::exit(0); }
    }

    if matches.opt_present("help") {
        print_usage(opts);
        unsafe { libc::exit(0); }
    }

    flag_map.insert("cookie".to_string(), {
        let ck = match matches.opt_str("cookie") {
            Some(a) => a,
            None => "".to_string()
        };
        let mut cal = Md5::new();
        cal.input_str(&ck);
        cal.result_str().to_string()
    });

    let ip = match matches.opt_str("ip") {
        Some(a) => a,
        None => "127.0.0.1".to_string()
    };

    flag_map.insert("addr".to_string(), {
        let port = match matches.opt_str("master-port") {
            Some(a) => a,
            None => "4400".to_string()
        };
        ip + ":" + &port
    });

    flag_map.insert("slave-port".to_string(), {
        match matches.opt_str("slave-port") {
            Some(a) => a,
            None => "4404".to_string()
        }
    });

    flag_map
}

fn print_version() {
    println!("saira-client {}", env!("CARGO_PKG_VERSION"));
}

fn print_usage(opts: Options) {
    print!("{}", opts.usage("Usage: saira-slave [OPTIONS]"));
}

