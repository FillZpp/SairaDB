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


use std::io::{stderr, Write};
use super::libc;


pub enum State<T> {
    Half(T),
    Done(T)
}

pub enum Operations {
    None,

    ShowDBs,
    Create(String),
    Drop(String),
    Use(String),
    
    Get(String, State<Vec<String>>),
    Set(String, String, i32),
    Add(String, String, i32),
    Del(String, State<Vec<String>>)
}

#[derive(Debug)]
#[allow(non_snake_case)]
#[derive(RustcDecodable, RustcEncodable)]
pub struct Query {
    Operation: String,
    Name: String,
    Attributes: Vec<String>,
    Data: String
}

impl Query {
    pub fn new(oper: Operations) -> Query {
        match oper {
            Operations::None =>
                Query {
                    Operation: "none".to_string(),
                    Name: "".to_string(),
                    Attributes: Vec::new(),
                    Data: "".to_string(),
                },

            Operations::ShowDBs =>
                Query {
                    Operation: "show_dbs".to_string(),
                    Name: "".to_string(),
                    Attributes: Vec::new(),
                    Data: "".to_string(),
                },

            Operations::Create(name) =>
                Query {
                    Operation: "create".to_string(),
                    Name: name,
                    Attributes: Vec::new(),
                    Data: "".to_string(),
                },

            Operations::Drop(name) =>
                Query {
                    Operation: "drop".to_string(),
                    Name: name,
                    Attributes: Vec::new(),
                    Data: "".to_string(),
                },

            Operations::Use(name) =>
                Query {
                    Operation: "use".to_string(),
                    Name: name,
                    Attributes: Vec::new(),
                    Data: "".to_string(),
                },

            Operations::Get(key, attrs) => {
                let attrs = match attrs {
                    State::Done(a) => a,
                    State::Half(_) => {
                        let _ = writeln!(stderr(), "Error: operations error");
                        unsafe { libc::exit(4); }
                    }
                };
                Query {
                    Operation: "get".to_string(),
                    Name: key,
                    Attributes: attrs,
                    Data: "".to_string()
                }
            }

            Operations::Set(key, data, _) =>
                Query {
                    Operation: "set".to_string(),
                    Name: key,
                    Attributes: Vec::new(),
                    Data: data
                },

            Operations::Add(key, data, _) =>
                Query {
                    Operation: "add".to_string(),
                    Name: key,
                    Attributes: Vec::new(),
                    Data: data
                },

            Operations::Del(key, attrs) => {
                let attrs = match attrs {
                    State::Done(a) => a,
                    State::Half(_) => {
                        let _ = writeln!(stderr(), "Error: operations error");
                        unsafe { libc::exit(4); }
                    }
                };
                Query {
                    Operation: "del".to_string(),
                    Name: key,
                    Attributes: attrs,
                    Data: "".to_string()
                }
            }
        }
    }
}


