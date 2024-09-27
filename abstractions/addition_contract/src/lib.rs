use serde::{Deserialize, Serialize};
use wasm_macros::wasm_export;

// Bring in alloc and dealloc functions
use wasm_lib::alloc;

#[derive(Serialize, Deserialize)]
pub struct AddThreeNumsReq {
    pub a: u32,
    pub b: u32,
    pub c: u32,
}

#[derive(Serialize, Deserialize)]
pub struct AddFourNumsReq {
    pub a: u32,
    pub b: u32,
    pub c: u32,
    pub d: u32
}


#[derive(Serialize, Deserialize)]
pub struct JoinTwoStrings {
    pub a: String,
    pub b: String,
}


#[wasm_export]
pub fn add_three_nums(input: AddThreeNumsReq) -> String {
    let sum = input.a + input.b + input.c;
    sum.to_string()
}

#[wasm_export]
pub fn add_four_nums(input: AddFourNumsReq) -> String {
    let sum = input.a + input.b + input.c + input.d;
    sum.to_string()
}

#[wasm_export]
fn concatenate_strings(input: JoinTwoStrings) -> String {
    format!("{}-{}", input.a, input.b)
}

