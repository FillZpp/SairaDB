// SairaDB - A distributed database
// Copyright (C) 2015 by Siyu Wang
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.


#![feature(libc)]

extern crate libc;

mod sr_args;
mod sr_readline;
mod sr_query;
mod sr_conn;

use std::net::TcpStream;
use std::io::{stdout, stderr, Write};


fn main() {
    let flag_map = sr_args::get_flags();

    let stream = {
        let addr: &str = flag_map.get("addr").unwrap();
        match TcpStream::connect(addr) {
            Ok(s) => s,
            Err(e) => {
                let _ = writeln!(stderr(), "Error: Can not connect to {}\n{}",
                                 addr, e);
                unsafe { libc::exit(4); }
            }
        }
    };
    
    let _ = writeln!(stdout(), "SairaDB Client {}", env!("CARGO_PKG_VERSION"));
    sr_conn::start_repl(stream, flag_map);
}

