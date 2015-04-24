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


use std::ffi::{CStr, CString};
use std::str::from_utf8;
use std::io::stdin;

#[cfg(unix)]
#[link(name = "readline")]
extern "C" {
    #[link_name = "readline"]
    fn rl_readline(prompt: *const i8) -> *const i8;
}

#[cfg(unix)]
fn read_line_unix(prompt: &str) -> Option<String> {
    let pr = CString::new(prompt.as_bytes()).unwrap();
    let sp = unsafe { rl_readline(pr.as_ptr()) };

    if sp.is_null() {
        None
    } else {
        let cs = unsafe { CStr::from_ptr(sp) };
        Some(from_utf8(cs.to_bytes()).unwrap().to_string())
    }
}

pub fn read_line(prompt: &str) -> Option<String> {
    if cfg!(unix) {
        read_line_unix(prompt)
    } else {
        let mut s = "".to_string();
        print!("{}", prompt);
        match stdin().read_line(&mut s) {
            Ok(_) => {
                s.pop();  // pop '\n'
                Some(s)
            }
            Err(_) => None
        }
    }
}



