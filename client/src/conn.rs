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

extern crate rustc_unicode;

use std::net::TcpStream;
use std::collections::HashMap;
use std::io::{stderr, Write, Read};
use std::str;
use self::rustc_unicode::str::UnicodeStr;
use super::libc;
use super::rustc_serialize::*;
use super::readline;
use super::query::*;


fn do_write(stream: &mut TcpStream, msg: &String) {
    match stream.write_all(msg.as_bytes()) {
        Ok(_) => {}
        Err(e) => {
            let _ = writeln!(stderr(), "Error: {}", e);
            unsafe { libc::exit(4); }
        }
    }
}

fn do_read(stream: &mut TcpStream, buf: &mut [u8]) -> String {
    let mut msg = "".to_string();
    loop {
        let n = match stream.read(buf) {
            Ok(n) => n,
            Err(e) => {
                let _ = writeln!(stderr(), "Error: {}", e);
                unsafe { libc::exit(4); }
            }
        };
        msg = msg + match str::from_utf8(&buf[0..n]) {
            Ok(s) => s,
            Err(e) => {
                let _ = writeln!(stderr(), "Error: {}", e);
                unsafe { libc::exit(4); }
            }
        };
        if n < buf.len() {
            break;
        }
    }
    msg
}

fn read_line(prompt: &str) -> String {
    match readline::read_line(prompt) {
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

fn check_data(data: &str) -> i32 {
    let mut ret = 0;
    let mut ins = false;
    let mut fxx = false;

    for c in data.chars() {
        if ins {
            if fxx {
                fxx = false;
                continue;
            }
            match c {
                '\"' => ins = false,
                '\\' => fxx = true,
                _ => {}
            }
        } else {
            match c {
                '\"' => ins = true,
                '{' => ret += 1,
                '}' => ret -= 1,
                _ => {}
            }
        }
    }
    ret
}

fn print_help() {
    println!("\nHelp:");
    println!("\"help\" to print commands");
    println!("\"quit\" to exit");
    println!("");
    println!("\"show dbs\"         print a list of databases");
    println!("\"create <db_name>\" create new database");
    println!("\"drop <db_name>\"   drop an exist database");
    println!("\"use <db_name>\"    focus on an exist database");
    println!("");
    println!("\"get <key> [attributes]\"");
    println!("\"set <key> <json>\"");
    println!("\"add <key> <json>\"");
    println!("\"del <key> [attributes]\"\n");
}

pub fn start_repl(flag_map: HashMap<String, String>) {
    let mut stream;
    let mut addr: String = flag_map.get("addr").unwrap().to_string();
    let slave_port: String = flag_map.get("slave-port").unwrap().to_string();
    let mut buf = [0u8; 100];
    let mut cache: HashMap<String, String> = HashMap::new();
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
        let msg: Vec<String> = do_decode(&do_read(&mut stream, &mut buf));
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
    let mut cur_db = "default".to_string();
    loop {
        let mut line = "".to_string();
        let to_readline = match &operation {
            &Operations::Get(_, State::Done(_)) => false,
            &Operations::Del(_, State::Done(_)) => false,
            &Operations::Set(_, _, n) => if n <= 0 {false} else {true},
            &Operations::Add(_, _, n) => if n <= 0 {false} else {true},
            _ => true
        };
        if to_readline {
            loop {
                line = match operation {
                    Operations::None => read_line("saira>> "),
                    _ => read_line("saira.> ")
                };
                if line != "".to_string() {
                    break;
                }
            }
        }

        match operation {
            Operations::None => {
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
                            println!("(error) wrong use of command.");
                            println!("Type 'help' to get a help list.");
                            continue;
                        }

                        let qry = vec!["show_dbs".to_string()];
                        do_write(&mut stream, &do_encode(&qry));
                        let dbs: Vec<String> = do_decode(&do_read(&mut stream,
                                                                  &mut buf));
                        for db in dbs {
                            println!("{}", db);
                        }
                    }

                    "create" => {
                        let name = match words.next() {
                            Some(n) => n.to_string(),
                            None => continue
                        };
                        
                        let qry = vec!["create".to_string(), name];
                        do_write(&mut stream, &do_encode(&qry));
                        let res = do_read(&mut stream, &mut buf);
                        println!("{}", res);
                    }

                    "drop" => {
                        let name = match words.next() {
                            Some(n) => n.to_string(),
                            None => continue
                        };

                        let qry = vec!["drop".to_string(), name];
                        do_write(&mut stream, &do_encode(&qry));
                        let res = do_read(&mut stream, &mut buf);
                        println!("{}", res);
                    }

                    "use" => {
                        let name = match words.next() {
                            Some(n) => n.to_string(),
                            None => continue
                        };

                        let qry = vec!["use".to_string(), name.to_string()];
                        do_write(&mut stream, &do_encode(&qry));
                        let res = do_read(&mut stream, &mut buf);
                        if &res == "ok" {
                            cur_db = name;
                        }
                        println!("{}", res);
                    }

                    "get" => {
                        let mut rests = match words.next() {
                            Some(r) => r.trim(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        }.splitn(2, " ");

                        let key = match rests.next() {
                            Some(k) => k.to_string(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        };

                        operation = Operations::Get(key, {
                            if let Some(rest) = rests.next() {
                                let mut attrs = Vec::new();
                                let rest = rest.trim();

                                if !rest.starts_with("(") {
                                    println!("(error) wrong use of command.");
                                    println!("Type 'help' to get a help list.");
                                    continue;
                                }
                                
                                let words: Vec<&str> = rest[1..].split(',')
                                    .collect();
                                let mut check = true;
                                for i in 0..(words.len()-1) {
                                    let s = words[i].trim();
                                    if s == "" || !s.is_alphanumeric() {
                                        println!("(error) wrong use of command.");
                                        println!("Type 'help' to get a help list.");
                                        check = false;
                                        break;
                                    }
                                    attrs.push(s.to_string());
                                }
                                if !check {
                                    continue;
                                }

                                let last = words[words.len() - 1];
                                if last == "" {
                                    State::Half(attrs)
                                } else if last.ends_with(")") {
                                    let s = &last[..(last.len() - 1)].trim();
                                    if !s.is_alphanumeric() {
                                        println!("(error) wrong use of command.");
                                        println!("Type 'help' to get a help list.");
                                        continue;
                                    }
                                    if s != &"" {
                                        attrs.push(s.to_string());
                                    }
                                    State::Done(attrs)
                                } else {
                                    println!("(error) parsing error.");
                                    println!("Type 'help' to get a help list.");
                                    continue;
                                }
                            } else {
                                State::Done(Vec::new())
                            }
                        });
                    }

                    "set" => {
                        let mut rests = match words.next() {
                            Some(r) => r.trim(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        }.splitn(2, " ");

                        let key = match rests.next() {
                            Some(k) => k.to_string(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        };

                        let data = match rests.next() {
                            Some(d) => d.trim(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        };
                        if !data.starts_with("{") {
                            println!("(error) wrong use of command.");
                            println!("Type 'help' to get a help list.");
                            continue;
                        }
                        let n = check_data(&data);
                        operation = Operations::Set(key, data.to_string(), n);
                    }

                    "add" => {
                        let mut rests = match words.next() {
                            Some(r) => r.trim(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        }.splitn(2, " ");

                        let key = match rests.next() {
                            Some(k) => k.to_string(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        };

                        let data = match rests.next() {
                            Some(d) => d.trim(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        };
                        if !data.starts_with("{") {
                            println!("(error) wrong use of command.");
                            println!("Type 'help' to get a help list.");
                            continue;
                        }
                        let n = check_data(&data);
                        operation = Operations::Add(key, data.to_string(), n);
                    }

                    "del" => {
                        let mut rests = match words.next() {
                            Some(r) => r.trim(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        }.splitn(2, " ");

                        let key = match rests.next() {
                            Some(k) => k.to_string(),
                            None => {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                        };

                        operation = Operations::Del(key, {
                            if let Some(rest) = rests.next() {
                                let mut attrs = Vec::new();
                                let rest = rest.trim();

                                if !rest.starts_with("(") {
                                    println!("(error) wrong use of command.");
                                    println!("Type 'help' to get a help list.");
                                    continue;
                                }
                                
                                let words: Vec<&str> = rest[1..].split(',')
                                    .collect();
                                let mut check = true;
                                for i in 0..(words.len()-1) {
                                    let s = words[i].trim();
                                    if s == "" || !s.is_alphanumeric() {
                                        println!("(error) wrong use of command.");
                                        println!("Type 'help' to get a help list.");
                                        check = false;
                                        break;
                                    }
                                    attrs.push(s.to_string());
                                }
                                if !check {
                                    continue;
                                }

                                let last = words[words.len() - 1];
                                if last == "" {
                                    State::Half(attrs)
                                } else if last.ends_with(")") {
                                    let s = &last[..(last.len() - 1)].trim();
                                    if !s.is_alphanumeric() {
                                        println!("(error) wrong use of command.");
                                        println!("Type 'help' to get a help list.");
                                        continue;
                                    }
                                    if s != &"" {
                                        attrs.push(s.to_string());
                                    }
                                    State::Done(attrs)
                                } else {
                                    println!("(error) parsing error.");
                                    println!("Type 'help' to get a help list.");
                                    continue;
                                }
                            } else {
                                State::Done(Vec::new())
                            }
                        });
                    }
                        
                    other => println!("Error: unknown command '{}'", other),
                } // match cmd
            }

            Operations::Get(key, attrs) => {
                operation = Operations::None;
                let attrs = match attrs {
                    State::Half(mut attrs) => {
                        let words: Vec<&str> = line.split(",").collect();
                        let mut check = true;
                        for i in 0..(words.len()-1) {
                            let s = words[i].trim();
                            if s == "" || !s.is_alphanumeric() {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                check = false;
                                break;
                            }
                            attrs.push(s.to_string());
                        }
                        if !check {
                            continue;
                        }

                        let last = words[words.len() - 1];
                        if last == "" {
                            operation = Operations::Get(key, State::Half(attrs));
                            continue;
                        } else if last.ends_with(")") {
                            let s = &last[..(last.len() - 1)].trim();
                            if !s.is_alphanumeric() {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                            if s != &"" {
                                attrs.push(s.to_string());
                            }
                            attrs
                        } else {
                            println!("(error) parsing error.");
                            println!("Type 'help' to get a help list.");
                            continue;
                        }
                    }
                    State::Done(attrs) => attrs,
                };
                
                let qry = do_encode(
                    &Query::new(cur_db.to_string(),
                        Operations::Get(key.to_string(), State::Done(attrs))));

                match cache.get(&key) {
                    Some(ip) => {
                        let addr: &str = &(ip.to_string() + &slave_port);
                        match TcpStream::connect(addr) {
                            Ok(mut stream) => {
                                do_write(&mut stream, &qry);
                                let res: Vec<String> = do_decode(
                                    &do_read(&mut stream, &mut buf));
                                // TODO
                                // Check the result
                                // If get nothing, then ask master for slave ip
                                println!("{:?}", res);
                                continue;
                            }
                            Err(_) => {}
                        }
                    }
                    None => {}
                }

                // ask master for slave ip
                do_write(&mut stream, &do_encode(&vec!["key", &key]));
                let res: Vec<String> = do_decode(&do_read(&mut stream, &mut buf));
                match res[0].as_ref() {
                    "ok" => {
                        cache.insert(key, res[1].to_string());
                        let addr: &str = &(res[1].to_string() + &slave_port);
                        match TcpStream::connect(addr) {
                            Ok(mut stream) => {
                                do_write(&mut stream, &qry);
                                let res: Vec<String> = do_decode(
                                    &do_read(&mut stream, &mut buf));
                                println!("{:?}", res);
                            }
                            Err(e) => println!("(error) can not connect to slave: {}", e)
                        }
                    }
                    _ => println!("{}", res[1])
                }
            }

            Operations::Set(key, data, mut n) => {
                operation = Operations::None;
                if n < 0 {
                    println!("(error) parsing error.");
                    println!("Type 'help' to get a help list.");
                    continue;
                } else if n > 0 {
                    n += check_data(&line);
                    operation = Operations::Set(key, data + &line, n);
                    continue;
                }

                if let Err(_) = json::Json::from_str(&data) {
                    println!("(error) parsing error.");
                    println!("Type 'help' to get a help list.");
                    continue;
                }
                
                let qry = &do_encode(&Query::new(cur_db.to_string(),
                    Operations::Set(key.to_string(), data, n)));

                match cache.get(&key) {
                    Some(ip) => {
                        let addr: &str = &(ip.to_string() + &slave_port);
                        match TcpStream::connect(addr) {
                            Ok(mut stream) => {
                                do_write(&mut stream, &qry);
                                let res: Vec<String> = do_decode(
                                    &do_read(&mut stream, &mut buf));
                                // TODO
                                // Check the result
                                // If get nothing, then ask master for slave ip
                                println!("{:?}", res);
                                continue;
                            }
                            Err(_) => {}
                        }
                    }
                    None => {}
                }

                // ask master for slave ip
                do_write(&mut stream, &do_encode(&vec!["key", &key]));
                let res: Vec<String> = do_decode(&do_read(&mut stream, &mut buf));
                match res[0].as_ref() {
                    "ok" => {
                        cache.insert(key, res[1].to_string());
                        let addr: &str = &(res[1].to_string() + &slave_port);
                        match TcpStream::connect(addr) {
                            Ok(mut stream) => {
                                do_write(&mut stream, &qry);
                                let res: Vec<String> = do_decode(
                                    &do_read(&mut stream, &mut buf));
                                println!("{:?}", res);
                            }
                            Err(e) => println!("(error) can not connect to slave: {}", e)
                        }
                    }
                    _ => println!("{}", res[1])
                }
            }

            Operations::Add(key, data, mut n) => {
                operation = Operations::None;
                if n < 0 {
                    println!("(error) parsing error.");
                    println!("Type 'help' to get a help list.");
                    continue;
                } else if n > 0 {
                    n += check_data(&line);
                    operation = Operations::Add(key, data + &line, n);
                    continue;
                }

                if let Err(_) = json::Json::from_str(&data) {
                    println!("(error) parsing error.");
                    println!("Type 'help' to get a help list.");
                    continue;
                }
                
                let qry = &do_encode(&Query::new(cur_db.to_string(),
                    Operations::Add(key.to_string(), data, n)));

                match cache.get(&key) {
                    Some(ip) => {
                        let addr: &str = &(ip.to_string() + &slave_port);
                        match TcpStream::connect(addr) {
                            Ok(mut stream) => {
                                do_write(&mut stream, &qry);
                                let res: Vec<String> = do_decode(
                                    &do_read(&mut stream, &mut buf));
                                // TODO
                                // Check the result
                                // If get nothing, then ask master for slave ip
                                println!("{:?}", res);
                                continue;
                            }
                            Err(_) => {}
                        }
                    }
                    None => {}
                }

                // ask master for slave ip
                do_write(&mut stream, &do_encode(&vec!["key", &key]));
                let res: Vec<String> = do_decode(&do_read(&mut stream, &mut buf));
                match res[0].as_ref() {
                    "ok" => {
                        cache.insert(key, res[1].to_string());
                        let addr: &str = &(res[1].to_string() + &slave_port);
                        match TcpStream::connect(addr) {
                            Ok(mut stream) => {
                                do_write(&mut stream, &qry);
                                let res: Vec<String> = do_decode(
                                    &do_read(&mut stream, &mut buf));
                                println!("{:?}", res);
                            }
                            Err(e) => println!("(error) can not connect to slave: {}", e)
                        }
                    }
                    _ => println!("{}", res[1])
                }
            }

            Operations::Del(key, attrs) => {
                operation = Operations::None;
                let attrs = match attrs {
                    State::Half(mut attrs) => {
                        let words: Vec<&str> = line.split(",").collect();
                        let mut check = true;
                        for i in 0..(words.len()-1) {
                            let s = words[i].trim();
                            if s == "" || !s.is_alphanumeric() {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                check = false;
                                break;
                            }
                            attrs.push(s.to_string());
                        }
                        if !check {
                            continue;
                        }

                        let last = words[words.len() - 1];
                        if last == "" {
                            operation = Operations::Del(key, State::Half(attrs));
                            continue;
                        } else if last.ends_with(")") {
                            let s = &last[..(last.len() - 1)].trim();
                            if !s.is_alphanumeric() {
                                println!("(error) wrong use of command.");
                                println!("Type 'help' to get a help list.");
                                continue;
                            }
                            if s != &"" {
                                attrs.push(s.to_string());
                            }
                            attrs
                        } else {
                            println!("(error) parsing error.");
                            println!("Type 'help' to get a help list.");
                            continue;
                        }
                    }
                    State::Done(attrs) => attrs,
                };

                let qry = &do_encode(&Query::new(cur_db.to_string(),
                    Operations::Del(key.to_string(), State::Done(attrs))));

                match cache.get(&key) {
                    Some(ip) => {
                        let addr: &str = &(ip.to_string() + &slave_port);
                        match TcpStream::connect(addr) {
                            Ok(mut stream) => {
                                do_write(&mut stream, &qry);
                                let res: Vec<String> = do_decode(
                                    &do_read(&mut stream, &mut buf));
                                // TODO
                                // Check the result
                                // If get nothing, then ask master for slave ip
                                println!("{:?}", res);
                                continue;
                            }
                            Err(_) => {}
                        }
                    }
                    None => {}
                }

                // ask master for slave ip
                do_write(&mut stream, &do_encode(&vec!["key", &key]));
                let res: Vec<String> = do_decode(&do_read(&mut stream, &mut buf));
                match res[0].as_ref() {
                    "ok" => {
                        cache.insert(key, res[1].to_string());
                        let addr: &str = &(res[1].to_string() + &slave_port);
                        match TcpStream::connect(addr) {
                            Ok(mut stream) => {
                                do_write(&mut stream, &qry);
                                let res: Vec<String> = do_decode(
                                    &do_read(&mut stream, &mut buf));
                                println!("{:?}", res);
                            }
                            Err(e) => println!("(error) can not connect to slave: {}", e)
                        }
                    }
                    _ => println!("{}", res[1])
                }
            }
        } // match operation
    }
}

