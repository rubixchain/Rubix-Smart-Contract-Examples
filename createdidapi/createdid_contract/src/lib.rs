use std::mem;
use std::ffi::{CString, CStr};
use std::os::raw::c_void;
extern crate serde;
extern crate serde_json;
#[macro_use] extern crate serde_derive;

#[derive(Serialize, Deserialize, Debug)]
struct SCTDataReply {
    BlockNo: u32,
    BlockId: String,
    SmartContractData: String,
}

#[derive(Serialize, Deserialize, Debug,Clone)]
struct SmartContractData {
    Type: i32,
    PrivPWD: String,
    ImgFile:String,
}


/*This alloc function is used to allocate 1024 bytes and returns a pointer.
When this function is called in the Go code we will receive the pointer.
Whatever data we need is being pushed onto this memory location.
*/
#[no_mangle]
pub extern "C" fn alloc() -> *mut c_void {
    let mut buf = Vec::with_capacity(1024);
    let ptr = buf.as_mut_ptr();

    mem::forget(buf);

    ptr
}
extern "C" {
   
    fn create_did(ptr: *const u8);
}
#[no_mangle]
pub unsafe extern "C" fn dealloc(ptr: *mut c_void) {
    let _ = Vec::from_raw_parts(ptr, 0, 1024);
}

// Extractig smartcontract data from the token chain data
fn extract_smartcontract_data(blocks: &[SCTDataReply]) -> Vec<SmartContractData> {   
    let mut vec_sc_data: Vec<SmartContractData> = Vec::new();


    for block in blocks {
        let  scdata = &block.SmartContractData;
        if scdata.is_empty() {
            continue;
        }
        if let Ok(data) = serde_json::from_str::<SmartContractData>(&block.SmartContractData) {
           vec_sc_data.push(data);
        }
    }

    vec_sc_data
}
/* This is the did function which is being triggered from Go
Here we are taking the entire tokenchain Data and we are calling the createDid Api from the contract  */
#[no_mangle]
pub unsafe fn did(ptr: *mut u8) {
    
     let json_data = CStr::from_ptr(ptr as *const i8).to_str().unwrap();
    // Deserialize the JSON data into a vector of SCTDataReply structs
    let mut blocks: Vec<SCTDataReply> = serde_json::from_str(json_data).expect("Failed to deserialize JSON");
 
    let smartcontract_data_vec = extract_smartcontract_data(&blocks);
    let vec_len = smartcontract_data_vec.len();
    let last_element =  smartcontract_data_vec[vec_len-1].clone();
    // let first_element_ptr: *const SmartContractData = first_element as *const SmartContractData;
    let mut serialized = serde_json::to_string(&last_element).unwrap();
    // if 1==1{rbt_transfer()}
    // let length_serialized_data = serialized.len();
    create_did(serialized.as_mut_ptr());
    
   
}