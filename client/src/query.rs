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


pub enum Operations {
    None,

    ShowDBs,
    CreateDB(String, String),
    DropDB(String),
    Use(String),
    
    Select(Option<Vec<String>>, Option<Vec<String>>),
    Insert(Vec<String>),
    Update(Option<Vec<String>>, Option<Vec<String>>),
    Delete(Option<Vec<String>>, Option<Vec<String>>)
}

#[derive(RustcDecodable, RustcEncodable)]
pub struct Query {
    operation: String,
    name: String,
    attributes: Vec<String>,
    data: Vec<String>,
    conditions: Vec<String>,
}

impl Query {
    pub fn new(oper: Operations) -> Query {
        match oper {
            Operations::None =>
                Query {
                    operation: "none".to_string(),
                    name: "".to_string(),
                    attributes: Vec::new(),
                    data: Vec::new(),
                    conditions: Vec::new(),
                },

            Operations::ShowDBs =>
                Query {
                    operation: "show_dbs".to_string(),
                    name: "".to_string(),
                    attributes: Vec::new(),
                    data: Vec::new(),
                    conditions: Vec::new(),
                },

            Operations::CreateDB(name, key) =>
                Query {
                    operation: "create_db".to_string(),
                    name: name,
                    attributes: vec![key],
                    data: Vec::new(),
                    conditions: Vec::new(),
                },

            Operations::DropDB(name) =>
                Query {
                    operation: "drop_db".to_string(),
                    name: name,
                    attributes: Vec::new(),
                    data: Vec::new(),
                    conditions: Vec::new(),
                },

            Operations::Use(name) =>
                Query {
                    operation: "use".to_string(),
                    name: name,
                    attributes: Vec::new(),
                    data: Vec::new(),
                    conditions: Vec::new(),
                },

            Operations::Select(attrs, conds) => {
                let attrs = match attrs {
                    None => Vec::new(),
                    Some(a) => a
                };

                let conds = match conds {
                    None => Vec::new(),
                    Some(a) => a
                };
                
                Query {
                    operation: "select".to_string(),
                    name: "".to_string(),
                    attributes: attrs,
                    data: Vec::new(),
                    conditions: conds,
                }
            }

            Operations::Insert(data) =>
                Query {
                    operation: "insert".to_string(),
                    name: "".to_string(),
                    attributes: Vec::new(),
                    data: data,
                    conditions: Vec::new(),
                },

            Operations::Update(data, conds) => {
                let data = match data {
                    None => Vec::new(),
                    Some(a) => a
                };

                let conds = match conds {
                    None => Vec::new(),
                    Some(a) => a
                };
                
                Query {
                    operation: "update".to_string(),
                    name: "".to_string(),
                    attributes: Vec::new(),
                    data: data,
                    conditions: conds,
                }
            }

            Operations::Delete(attrs, conds) => {
                let attrs = match attrs {
                    None => Vec::new(),
                    Some(a) => a
                };

                let conds = match conds {
                    None => Vec::new(),
                    Some(a) => a
                };
                
                Query {
                    operation: "delete".to_string(),
                    name: "".to_string(),
                    attributes: attrs,
                    data: Vec::new(),
                    conditions: conds,
                }
            }
        }
    }
}


