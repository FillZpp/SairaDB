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


extern crate serialize;

use std::collections::{BTreeMap, HashMap, VecDeque};
use std::sync::atomic::AtomicUsize;
//use std::sync::Arc;
//use std::sync::mpsc::{Sender, Receiver};
use self::serialize::json::Json;
//use super::sr_type::{Types, BasicTypes};


#[allow(dead_code)]
struct Page {
    id: u64,
    size: AtomicUsize,
    units: BTreeMap<String, Json>
}

#[allow(dead_code)]
struct Set {
    name: String,
    key: String,
    size: AtomicUsize,
    pages: VecDeque<Page>,
}

#[allow(dead_code)]
pub struct Database {
    sets: HashMap<String, Set>,
}



