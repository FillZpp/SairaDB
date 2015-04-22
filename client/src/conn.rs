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


use std::net::TcpStream;
use std::collections::HashMap;
use std::io::{stdout, stderr, Write, Read};
use super::libc;
use super::rustc_serialize::*;
use super::readline;
use super::query::{Operations, Query};


fn do_write(stream: &mut TcpStream, msg: &String) {
    match stream.write_all(msg.as_bytes()) {
        Ok(_) => {}
        Err(e) => {
            let _ = writeln!(stderr(), "Error: {}", e);
            unsafe { libc::exit(4); }
        }
    }
}

fn do_read(stream: &mut TcpStream) -> String {
    let mut buf = [0u8; 32];
    let n = match stream.read(&mut buf) {
        Ok(n) => n,
        Err(e) => {
            let _ = writeln!(stderr(), "Error: {}", e);
            unsafe { libc::exit(4); }
        }
    };
    String::from_utf8_lossy(&buf[0..n]).into_owned()
}

fn read_line() -> String {
    match readline::read_line() {
        Some(s) => s,
        None => {
            unsafe { libc::exit(0); }
        }
    }
}

fn do_encode<T: Encodable>(object: &T) -> String {
    match json::encode(object) {
        Ok(s) => s,
        Err(_) => {
            let _ = writeln!(stderr(), "Error: json encode error");
            unsafe { libc::exit(0); }
        }
    }
}

fn print_help() {
    println!("\nHelp:");
    println!("\"help\" to print commands");
    println!("\"quit\" to exit");
    println!("");
    println!("\"show dbs\"                           print a list of databases");
    println!("\"create db <db_name> key <key_name>\" create new database");
    println!("\"drop db <db_name>\"                  drop an exist database");
    println!("\"use <db_name>\"                      focus on an exist database");
    println!("");
    println!("\"select <attributes> <conditions>\"");
    println!("\"insert <data>\"");
    println!("\"update <data> <conditions>\"");
    println!("\"delete <attributes> <conditions>\"\n");
}

pub fn start_repl(flag_map: HashMap<String, String>) {
    let mut stream;
    let mut addr: String = flag_map.get("addr").unwrap().to_string();
    loop {
        stream = {
            let s: &str = &addr;
            match TcpStream::connect(s) {
                Ok(s) => s,
                Err(e) => {
                    let _ = writeln!(stderr(), "Error: Can not connect to {}\n{}",
                                     addr, e);
                    unsafe { libc::exit(4); }
                }
            }
        };
        
        do_write(&mut stream, flag_map.get("cookie").unwrap());
        let msg: Vec<String> = match json::decode(&do_read(&mut stream)) {
            Ok(m) => m,
            Err(e) => {
                let _ = writeln!(stderr(), "Error: json parsing error\n{}", e);
                unsafe { libc::exit(4); }
            }
        };

        if msg.len() == 1 {
            if msg[0] == "ok".to_string() {
                break;
            } else if msg[0] == "wrong".to_string() {
                let _ = writeln!(stderr(), "Error: wrong cookie");
                unsafe { libc::exit(0); }
            }
        } else if msg.len() == 2 && msg[0] == "redirect" {
            addr = msg[1].to_string();
            continue;
        }

        let _ = writeln!(stderr(), "Error: undefine message");
        unsafe { libc::exit(4); }
        
    }

    let _ = writeln!(stdout(), "SairaDB Client {}", env!("CARGO_PKG_VERSION"));

    let mut operation = Operations::None;
    loop {
        let mut line;
        loop {
            line = read_line();
            if line != "".to_string() {
                break;
            }
        }

        match &mut operation {
            &mut Operations::None => {
                let mut words = line.trim().splitn(2, " ");
                let cmd = words.next().unwrap();

                match cmd {
                    "quit" => break,
                    "help" => print_help(),
                    
                    other => {
                        let _ = writeln!(stderr(), "Error: unknown command '{}'",
                                         other);
                        continue;
                    }
                }
            }

            &mut Operations::Select(ref mut attrs, ref mut conds) => {
            }

            &mut Operations::Insert(ref mut data) => {
            }

            &mut Operations::Update(ref mut data, ref mut conds) => {
            }

            &mut Operations::Delete(ref mut data, ref mut conds) => {
            }

            _ => {}
        }
    }
}

