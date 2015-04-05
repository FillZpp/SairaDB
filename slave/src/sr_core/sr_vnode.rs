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


use std::collections::HashMap;
use std::sync::atomic::{AtomicUsize, ATOMIC_USIZE_INIT};
use super::sr_db::Database;


#[allow(dead_code)]
pub struct VNode {
    id: u64,
    dups: Vec<String>,
    master: AtomicUsize,
    dbs: HashMap<String, Database>
}

#[allow(dead_code)]
impl VNode {
    pub fn new(id: u64) -> VNode {
        VNode {
            id: id,
            dups: Vec::new(),
            master: ATOMIC_USIZE_INIT,
            dbs: HashMap::new()
        }
    }
}

