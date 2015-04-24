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


extern crate sys_info;

mod sr_db;
mod sr_vnode;

pub use self::sr_db::*;
pub use self::sr_vnode::*;
use std::sync::mpsc::Sender;

static mut td_num: u32 = 0;

pub fn init(log_sender: Sender<String>) {
    unsafe {
        td_num = match sys_info::cpu_num() {
            Ok(n) => if n < 2 { 2 } else { n },
            Err(e) => {
                let _ = log_sender.send("slave core get cpu num error: "
                                        .to_string() + &e);
                2
            }
        }
    }
}

