use serde::{Deserialize, Serialize};
use wasm_macros::wasm_export;

// Bring in alloc and dealloc functions
use wasm_lib::{alloc, dealloc};

#[derive(Serialize, Deserialize)]
pub struct AddThreeNumsReq {
    pub a: u32,
    pub b: u32,
    pub c: u32,
}

#[wasm_export]
pub fn add_three_nums(input: AddThreeNumsReq) -> String {
    let sum = input.a + input.b + input.c;
    sum.to_string()
}


