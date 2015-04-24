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
#[derive(RustcDecodable, RustcEncodable)]
pub struct Query {
    operation: String,
    name: String,
    attributes: Vec<String>,
    data: String
}

impl Query {
    pub fn new(oper: Operations) -> Query {
        match oper {
            Operations::None =>
                Query {
                    operation: "none".to_string(),
                    name: "".to_string(),
                    attributes: Vec::new(),
                    data: "".to_string(),
                },

            Operations::ShowDBs =>
                Query {
                    operation: "show_dbs".to_string(),
                    name: "".to_string(),
                    attributes: Vec::new(),
                    data: "".to_string(),
                },

            Operations::Create(name) =>
                Query {
                    operation: "create".to_string(),
                    name: name,
                    attributes: Vec::new(),
                    data: "".to_string(),
                },

            Operations::Drop(name) =>
                Query {
                    operation: "drop".to_string(),
                    name: name,
                    attributes: Vec::new(),
                    data: "".to_string(),
                },

            Operations::Use(name) =>
                Query {
                    operation: "use".to_string(),
                    name: name,
                    attributes: Vec::new(),
                    data: "".to_string(),
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
                    operation: "get".to_string(),
                    name: key,
                    attributes: attrs,
                    data: "".to_string()
                }
            }

            Operations::Set(key, data, _) =>
                Query {
                    operation: "set".to_string(),
                    name: key,
                    attributes: Vec::new(),
                    data: data
                },

            Operations::Add(key, data, _) =>
                Query {
                    operation: "add".to_string(),
                    name: key,
                    attributes: Vec::new(),
                    data: data
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
                    operation: "del".to_string(),
                    name: key,
                    attributes: attrs,
                    data: "".to_string()
                }
            }
        }
    }
}


