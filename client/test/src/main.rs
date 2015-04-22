extern crate rustc_serialize;

use std::net::{TcpListener, TcpStream, Shutdown};
use std::io::{Read, Write};
use std::str;
use std::thread;
use rustc_serialize::json;

#[derive(RustcDecodable, RustcEncodable)]
#[derive(Debug)]
struct Query {
    operation: String,
    name: String,
    attributes: Vec<String>,
    data: Vec<String>,
    conditions: Vec<String>
}

fn handler(mut stream: TcpStream) {
    println!("\nnew");
    let mut buf = [0u8; 100];
    stream.read(&mut buf).unwrap();
    println!("{}", str::from_utf8(&buf).unwrap());
    let _ = stream.write_all(b"[\"ok\"]");

    loop {
        let n = match stream.read(&mut buf) {
            Ok(n) => n,
            Err(_) => break
        };
        println!("{}", str::from_utf8(&buf[0..n]).unwrap());
        let qry: Query = match json::decode(str::from_utf8(&buf[0..n]).unwrap()) {
            Ok(q) => q,
            Err(e) => {
                println!("{}", e);
                break;
            }
        };

        if qry.operation == "quit".to_string() {
            break;
        }
        let _ = stream.write_all(b"ok");
    }
    let _ = stream.shutdown(Shutdown::Both);
    println!("end");
}

fn main() {
    let listener = TcpListener::bind("127.0.0.1:4400").unwrap();

    println!("start");
    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                thread::spawn(move ||{
                    handler(stream);
                });
            }
            Err(e) => {
                println!("{}", e);
            }
        }
    }
}

