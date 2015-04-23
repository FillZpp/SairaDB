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
use std::io::{stderr, Write, Read};
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
            let _ = writeln!(stderr(), "Error: json encode error.");
            unsafe { libc::exit(4); }
        }
    }
}

fn do_decode<T: Decodable>(s: &str) -> T {
    match json::decode(s) {
        Ok(t) => t,
        Err(_) => {
            let _ = writeln!(stderr(), "Error: json decode error.");
            unsafe { libc::exit(4); }
        }
    }
}

fn print_help() {
    println!("\nHelp:");
    println!("\"help\" to print commands");
    println!("\"quit\" to exit");
    println!("");
    println!("\"show dbs\"                           print a list of databases");
    println!("\"create <db_name> key <key_name>\" create new database");
    println!("\"drop <db_name>\"                  drop an exist database");
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
        let msg: Vec<String> = do_decode(&do_read(&mut stream));
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

    println!("SairaDB Client {}", env!("CARGO_PKG_VERSION"));

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

                    "show" => {
                        let mut check = false;
                        match words.next() {
                            Some(sec) => {
                                if sec.trim() == "dbs" {
                                    check = true;
                                }
                            }
                            
                            None => {}
                        }

                        if !check {
                            println!("Error: wrong command. Type 'help' to get a help list.");
                            continue;
                        }

                        let qry = Query::new(Operations::ShowDBs);
                        do_write(&mut stream, &do_encode(&qry));
                        let dbs: Vec<String> = do_decode(&do_read(&mut stream));
                        for db in dbs {
                            println!("{}", db);
                        }
                    }

                    "create" => {
                        let name = match words.next() {
                            Some(n) => n.to_string(),
                            None => continue
                        };
                        
                        let qry = Query::new(Operations::Create(name));
                        do_write(&mut stream, &do_encode(&qry));
                        let res = do_read(&mut stream);
                        println!("{}", res);
                    }

                    "drop" => {
                        let name = match words.next() {
                            Some(n) => n.to_string(),
                            None => continue
                        };

                        let qry = Query::new(Operations::Drop(name));
                        do_write(&mut stream, &do_encode(&qry));
                        let res = do_read(&mut stream);
                        println!("{}", res);
                    }

                    "use" => {
                        let name = match words.next() {
                            Some(n) => n.to_string(),
                            None => continue
                        };

                        let qry = Query::new(Operations::Use(name));
                        do_write(&mut stream, &do_encode(&qry));
                        let res = do_read(&mut stream);
                        println!("{}", res);
                    }

                    "select" => {
                        
                    }
                        
                    other => println!("Error: unknown command '{}'", other),
                } // match cmd
            }

            &mut Operations::Get(_, ref mut attrs) => {
            }

            &mut Operations::Set(_, ref mut data) => {
            }

            &mut Operations::Add(_, ref mut data) => {
            }

            &mut Operations::Del(_, ref mut attrs) => {
            }

            _ => {}
        }
    }
}

