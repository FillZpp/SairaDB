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


#![feature(libc)]
#![feature(path_ext)]
#![feature(lookup_host, ip_addr)]


extern crate libc;
extern crate crypto;

mod sr_prefix;
mod sr_args;
mod sr_config;
mod sr_log;
mod sr_core;
mod sr_masterctl;


fn main() {
    let conf_map = sr_args::get_flags();
    let masters = sr_config::init(&conf_map);

    let (log_sender, log_thread)
        = sr_log::init(conf_map.get("data-dir").unwrap().to_string());

    let vnodes = sr_core::init(conf_map.get("data-dir").unwrap().to_string(),
                               log_sender.clone());
    
    sr_masterctl::init(masters, vnodes, &conf_map, log_sender.clone());

    //std::thread::sleep_ms(10);
}



