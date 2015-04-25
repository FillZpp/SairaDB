extern crate rustc_serialize;

use std::net::{TcpListener, TcpStream, Shutdown};
use std::io::{Read, Write, Error};
use std::str;
use std::thread;
use rustc_serialize::json;

#[derive(RustcDecodable, RustcEncodable)]
#[derive(Debug)]
struct Query {
    operation: String,
    name: String,
    attributes: Vec<String>,
    data: String,
}

fn read(stream: &mut TcpStream, buf: &mut [u8]) -> Result<String, Error> {
    let mut msg = "".to_string();
    loop {
        let n = match stream.read(buf) {
            Ok(n) => n,
            Err(e) => return Err(e)
        };

        msg = msg + str::from_utf8(&buf[0..n]).unwrap();
        if n < buf.len() {
            break;
        }
    }
    Ok(msg)
}

fn handler(mut stream: TcpStream) {
    println!("\nnew");
    let mut buf = [0u8; 80];
    stream.read(&mut buf).unwrap();
    println!("{}", str::from_utf8(&buf).unwrap());
    let _ = stream.write_all(b"[\"ok\"]");

    loop {
        let msg = match read(&mut stream, &mut buf) {
            Ok(m) => m,
            Err(_) => break
        };
        println!("{}", msg);
        let qry: Query = match json::decode(&msg) {
            Ok(q) => q,
            Err(e) => {
                println!("{}", e);
                break;
            }
        };

        match qry.operation.as_ref() {
            "quit" => break,

            "show_dbs" => {
                let dbs = vec!["default", "test"];
                let _ = stream.write_all(json::encode(&dbs).unwrap().as_bytes());
            }

            _ => {
                let _ = stream.write_all(b"ok");
            }
        }

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

