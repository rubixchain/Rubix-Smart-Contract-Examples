use serde::{Deserialize, Serialize};
use wasm_macros::wasm_export;

// Bring in alloc and dealloc functions
use wasm_lib::alloc;

#[derive(Serialize, Deserialize)]
pub struct WhiteListDIDContractInput {
    pub did: String,
}

// Just a dummy check is present
#[wasm_export]
pub fn whitelist_did(input: WhiteListDIDContractInput) -> String {
    let input_did = input.did;
    
    // Check if DID provided is not empty
    if input_did.len() > 0 {
        input_did
    } else {
        "0".to_string()
    }
}

