// SairaDB - A distributed database
// Copyright (C) 2015 by Siyu Wang
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
//	This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.


use std::collections::{BTreeMap, HashMap};
use std::sync::Arc;
use std::sync::mpsc::{Sender, Receiver};
use super::sr_type::{Types, BasicTypes};


#[allow(dead_code)]
struct Unit {
    value: Types,
    attrs: Option<BTreeMap<String, Arc<Unit>>>,
}

#[allow(dead_code)]
struct Column {
    units: BTreeMap<BasicTypes, Unit>,
}

#[allow(dead_code)]
struct Alter {
    res: Receiver<String>,
    con: BTreeMap<String, Types>,
    alt: i32  // 0 for insert, 1 for delete con, 2 for delete all rows
}

#[allow(dead_code)]
struct Table {
    key: String,
    columns: HashMap<String, Column>,
    alter_sender: Sender<Alter>,
}

#[allow(dead_code)]
pub struct Database {
    tables: HashMap<String, Table>,
    alter_sender: Sender<Alter>
}

