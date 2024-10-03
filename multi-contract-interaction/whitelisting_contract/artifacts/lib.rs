use core::slice;
use std::{mem, ptr};
use serde_derive::{Deserialize, Serialize};

// input data struct
#[derive(Serialize, Deserialize)]
pub struct SmartContractData {
    pub did: String,
}

#[no_mangle]
pub extern "C" fn alloc(size: usize) -> *mut u8 {
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    mem::forget(buf); // Prevent Rust from freeing the memory
    ptr
}

#[no_mangle]
pub extern "C" fn dealloc(ptr: *mut u8, size: usize) {
    unsafe {
        let _ = Vec::from_raw_parts(ptr, size, size);
    }
}

// Input: {"did": "did123.."}
// Output: DID string 
#[no_mangle]
pub extern "C" fn white_list_did(input_str: *mut u8, input_len: usize, output_ptr_ptr: *mut *mut u8, output_len_ptr: *mut usize) {
    // Deserialise input data
    let smart_contract_data_bytes = unsafe {
        slice::from_raw_parts(input_str, input_len)
    };
    let smart_contract_data: SmartContractData = serde_json::from_slice(smart_contract_data_bytes).unwrap();
    
    // Check if the DID has been provided or not
    if smart_contract_data.did.len() > 0 {
        let serialised_output = serde_json::to_vec(&smart_contract_data.did).unwrap();

        let output_len = serialised_output.len();
        let output_ptr = alloc(output_len);
            
        unsafe {
            ptr::copy_nonoverlapping(serialised_output.as_ptr(), output_ptr, output_len);
            *output_ptr_ptr = output_ptr;
            *output_len_ptr = output_len;
        }
    }

}