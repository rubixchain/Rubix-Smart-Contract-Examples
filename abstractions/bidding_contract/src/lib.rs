use serde::{Deserialize, Serialize};
use wasm_macros::wasm_export;

// Bring in alloc and dealloc functions
use wasm_lib::alloc;

#[derive(Serialize, Deserialize, Clone)]
pub struct PlaceBidInput {
    pub did: String,
    pub bid: u32,
}

#[wasm_export]
pub fn place_bid(input: PlaceBidInput) -> String {
    let input_did = input.clone().did;
    let input_bid = input.clone().bid;

    // Check if DID provided is not empty
    if input_did.len() == 0 {
       return "0".to_string();
    }

    // Check if input bid is less than 0
    if input_bid <= 0 {
        return "0".to_string();
    }

    // Return the result AS/Is
    let result = input.clone();
    serde_json::to_string(&result).unwrap()
}

