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


use std::collections::{HashMap, HashSet, VecDeque, BTreeMap};
use std::cmp::{Ord, Ordering, PartialOrd, Eq, PartialEq};
use std::hash::{Hasher, Hash};


#[derive(Debug)]
pub enum BasicTypes {
    Nothing,
    
    Int8(i8),
    Int(i64),
    Uint8(u8),
    Uint(u64),

    String(String),
    Date([i32; 6]),
}

#[derive(Debug)]
pub enum Types {
    BasicType(BasicTypes),

    Tuple(i32, Vec<BasicTypes>),
    List(VecDeque<BasicTypes>),
    Map(HashMap<String, BasicTypes>),
    Set(HashSet<BasicTypes>),

    Json(BTreeMap<String, Types>),
    Any,
}

impl PartialEq for BasicTypes {
    fn eq(&self, other: &Self) -> bool {
        match self {
            &BasicTypes::Nothing => match other {
                &BasicTypes::Nothing => true,
                _ => false,
            },
            
            &BasicTypes::Int8(l) => match other {
                &BasicTypes::Int8(r) => l == r,
                _ => false,
            },

            &BasicTypes::Int(l) => match other {
                &BasicTypes::Int(r) => l == r,
                _ => false,
            },

            &BasicTypes::Uint8(l) => match other {
                &BasicTypes::Uint8(r) => l == r,
                _ => false,
            },

            &BasicTypes::Uint(l) => match other {
                &BasicTypes::Uint(r) => l == r,
                _ => false,
            },

            &BasicTypes::String(ref l) => match other {
                &BasicTypes::String(ref r) => l == r,
                _ => false,
            },

            &BasicTypes::Date(l) => match other {
                &BasicTypes::Date(r) => l == r,
                _ => false,
            },
        }
    }
}

impl Eq for BasicTypes {}

impl Hash for BasicTypes {
    fn hash<H: Hasher>(&self, state: &mut H) {
        match self {
            &BasicTypes::Nothing => 0.hash(state),
            &BasicTypes::Int8(a) => a.hash(state),
            &BasicTypes::Int(a) => a.hash(state),
            &BasicTypes::Uint8(a) => a.hash(state),
            &BasicTypes::Uint(a) => a.hash(state),
            &BasicTypes::String(ref a) => a.hash(state),
            &BasicTypes::Date(a) => a.hash(state),
        }
    }
}

impl PartialOrd for BasicTypes {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        match self {
            &BasicTypes::Nothing => match other {
                &BasicTypes::Nothing => Some(Ordering::Equal),
                _ => None,
            },
            
            &BasicTypes::Int8(l) => match other {
                &BasicTypes::Int8(r) => Some(l.cmp(&r)),
                _ => None,
            },

            &BasicTypes::Int(l) => match other {
                &BasicTypes::Int(r) => Some(l.cmp(&r)),
                _ => None,
            },

            &BasicTypes::Uint8(l) => match other {
                &BasicTypes::Uint8(r) => Some(l.cmp(&r)),
                _ => None,
            },

            &BasicTypes::Uint(l) => match other {
                &BasicTypes::Uint(r) => Some(l.cmp(&r)),
                _ => None,
            },

            &BasicTypes::String(ref l) => match other {
                &BasicTypes::String(ref r) => Some(l.cmp(&r)),
                _ => None,
            },

            &BasicTypes::Date(l) => match other {
                &BasicTypes::Date(r) => Some(l.cmp(&r)),
                _ => None,
            },
        }
    }
}

impl Ord for BasicTypes {
    fn cmp(&self, other: &Self) -> Ordering {
        match self {
            &BasicTypes::Nothing => match other {
                &BasicTypes::Nothing => Ordering::Equal,
                _ => Ordering::Less,
            },
            
            &BasicTypes::Int8(l) => match other {
                &BasicTypes::Int8(r) => l.cmp(&r),
                _ => Ordering::Less,
            },

            &BasicTypes::Int(l) => match other {
                &BasicTypes::Int(r) => l.cmp(&r),
                _ => Ordering::Less,
            },

            &BasicTypes::Uint8(l) => match other {
                &BasicTypes::Uint8(r) => l.cmp(&r),
                _ => Ordering::Less,
            },

            &BasicTypes::Uint(l) => match other {
                &BasicTypes::Uint(r) => l.cmp(&r),
                _ => Ordering::Less,
            },

            &BasicTypes::String(ref l) => match other {
                &BasicTypes::String(ref r) => l.cmp(&r),
                _ => Ordering::Less,
            },

            &BasicTypes::Date(l) => match other {
                &BasicTypes::Date(r) => l.cmp(&r),
                _ => Ordering::Less,
            },
        }
    }
}

