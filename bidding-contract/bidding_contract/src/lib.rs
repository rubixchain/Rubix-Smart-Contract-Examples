use std::mem;
use std::ffi::{CString, CStr};
use std::os::raw::c_void;
extern crate serde;
extern crate serde_json;
//use ecies::decrypt;
//pub mod decrypt;
//mod seal;
// pub mod seal;
#[macro_use] extern crate serde_derive;
#[derive(Serialize, Deserialize, Debug)]
struct SCTDataReply {
    BlockNo: u32,
    BlockId: String,
    SmartContractData:String,
}
//test comments .....10
#[derive(Serialize, Deserialize, Debug)]
struct SmartContractData {
    did: String,
    bid: f64,
}

static mut CONTRACT_DATA: SmartContractData = SmartContractData {
    did: String::new(),
    bid: 0.0,
};
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

#[no_mangle]
pub unsafe extern "C" fn dealloc(ptr: *mut c_void) {
    let _ = Vec::from_raw_parts(ptr, 0, 1024);
}
// fn decrypted_scdatareply(blocks: &mut [SCTDataReply]){
//     let encrypted_privkeypath = "/home/rubix/Sai-Rubix/rubixgoplatform/linux/node9/Rubix/TestNetDID/bafybmiahqapy3fvpqn4feoawdotnewu3zh4cpne54uivaceb2bpd2ihnja/pvtKey.pem";
//     for block in blocks{
//         let mut ciphertext = block.SmartContractData.clone().into_bytes();
//         let mut decryptedsmartcontractdata = decrypt::decrypt_smartcontract_data(encrypted_privkeypath, ciphertext);
//         block.SmartContractData = decryptedsmartcontractdata; 
//     }
// }

//The smart contract logic for finding the highest bidder
//fn find_highest_bid_did(blocks: &mut [SCTDataReply]) -> Option<(String, f64)> {
       //let hex_str = "042e8c8c71a791fc1a3f475041f63b2cea236f446d2255ba7d85f9295ff62b4f967d5458de63fa496ab9803c4db93ede81781160d9aaf1637eb19124a96c68b96be421b8f5de7384f4e8529f0d3f41e2605c1c449bf9d8b401fcf09aa975a29a2324515607ce794ce9c9362a957e2f431f152b6e76fb3b672d9c0fdbecffb5b13dc1ec9029c56a552aabf52adb54acf44ed3f5a3a8b592e2e50581190116c2d67d924630070ce0545296f7e337fe39534442";
   
    // for block in&mut *blocks{
    //     let mut ciphertext = block.SmartContractData.clone().into_bytes();
    //     let decrypted_key = [26, 74, 206, 110, 150, 148, 87, 32, 213, 102, 150, 120, 224, 105, 131, 103, 58, 95, 72, 72, 142, 240, 97, 25, 113, 39, 140, 138, 164, 82, 187, 147];
    //     println!("cipher text of smartcontract data is {:?}",ciphertext);
    //     let mut decryptedsmartcontractdata = decrypt::decrypt_smartcontract_data(&decrypted_key,ciphertext);
    //     block.SmartContractData = decryptedsmartcontractdata; 
    // }
    //let mut max_bid_info: Option<(String, f64)> = None;
fn find_highest_bid_did(blocks: &[SCTDataReply]) -> Option<(String, f64)> {
    let mut max_bid_info: Option<(String, f64)> = None;

    for block in blocks {
        if let Ok(data) = serde_json::from_str::<SmartContractData>(&block.SmartContractData) {
            match max_bid_info {
                Some((_, max_bid)) if data.bid > max_bid => {
                    max_bid_info = Some((data.did.clone(), data.bid));
                }
                None => {
                    max_bid_info = Some((data.did.clone(), data.bid));
                }
                _ => {}
            }
        }
    }

    max_bid_info
}
/* This is the bid function which is being triggered from Go
Here we are taking the entire tokenchain Data and is checking for the highest bid
Then we are returning the did and the highest bid amount */
#[no_mangle]
pub unsafe fn bid(ptr: *mut u8) {
   
    // Assume get_blocks() returns a valid JSON string pointer
    // For testing, we'll use the hardcoded JSON data directly
     let json_data = CStr::from_ptr(ptr as *const i8).to_str().unwrap();
    
    // let json_data = br#"
    // [
    //     {
    //         "BlockNo": 0,
    //         "BlockId": "0-434ba0614ddc0db1f4bb22b77591ea60a7c04f343aa3236c67841ea7d070f6c6",
    //         "SmartContractData": ""
    //     },
    //     {
    //         "BlockNo": 1,
    //         "BlockId": "1-2cdabd4f3c7d89e624afc94cbd07d05c8377cbb5a35efd1ad8a08793ec3f3b53",
    //         "SmartContractData": "{\"did\":\"bafybmiflgqlbcwedw2mtqvhxyyckb455fig5q6l6zfgvllgiii2zdgsma4\",\"bid\":40.01}"
    //     },
    //     {
    //         "BlockNo": 2,
    //         "BlockId": "2-a1a3eeeea4c12e73fb4f6c44432ddccc695eda1022f196a9150e3165f370373e",
    //         "SmartContractData": "{\"did\":\"bafybmiflgqlbcwedw2mtqvhxyyckb455fig5q6l6zfgvllgiii2zdgsma4\",\"bid\":50.01}"
    //     }
    // ]
    // "#;

    // Deserialize the JSON data into a vector of SCTDataReply structs
    //let mut  blocks: Vec<SCTDataReply> = serde_json::from_str(json_data).expect("Failed to deserialize JSON");
    let blocks: Vec<SCTDataReply> = serde_json::from_str(json_data).expect("Failed to deserialize JSON");

    //let encrypted_privkeypath = "/home/rubix/Sai-Rubix/rubixgoplatform/linux/node9/Rubix/TestNetDID/bafybmiahqapy3fvpqn4feoawdotnewu3zh4cpne54uivaceb2bpd2ihnja/pvtKey.pem";
    // for block in&mut *blocks{
    //     let mut encoded_ciphertext = block.SmartContractData.clone();
    //     if encoded_ciphertext.is_empty() {
    //         continue;
    //     }
    //     let decrypted_key = [26, 74, 206, 110, 150, 148, 87, 32, 213, 102, 150, 120, 224, 105, 131, 103, 58, 95, 72, 72, 142, 240, 97, 25, 113, 39, 140, 138, 164, 82, 187, 147];
    //     let decoded_ciphertext = hex::decode(encoded_ciphertext).expect("unable to decode ");
    //     println!("cipher text of smartcontract data is {:?}",decoded_ciphertext);
    //     let mut decryptedsmartcontractdata = decrypt_smartcontract_data(&decrypted_key,decoded_ciphertext);
    //     //let mut decryptedscdata_in_string = std::str::from_utf8(decryptedsmartcontractdata).expect("unable to convert decrypted sc data into string format");
    //     println!("decrypted text of smartcontract data is {:?}",decryptedsmartcontractdata);
    //     block.SmartContractData = decryptedsmartcontractdata.to_string(); 
    // }
    // let blocks_as_bytes = blocks.as_mut_slice();
    
    // Find the block with the highest bid
    // match find_highest_bid_did(&blocks) {
    //     Some((block_no, max_bid)) => println!("The block with the highest bid is BlockNo {} with a bid of {}", block_no, max_bid),
    //     None => println!("No valid bids found."),
    // }
    match find_highest_bid_did(&blocks) {
        Some((block_no, max_bid)) => {
            println!("The block with the highest bid is BlockNo {} with a bid of {}", block_no, max_bid);
            
            // Use unsafe block to modify the static mutable variable
            unsafe {
                CONTRACT_DATA.did = block_no.to_string(); // Assuming block_no can be converted to a string
                CONTRACT_DATA.bid = max_bid;
            }
        },
        None => println!("No valid bids found."),
    }
}

#[no_mangle]
pub unsafe extern "C" fn get_state(ptr: *mut u8) {
    // Export state as JSON
    let string_content = serde_json::to_string(&CONTRACT_DATA).unwrap();
    get_return_string(string_content, ptr);
}

#[no_mangle]
unsafe fn get_return_string(string_content: String, ptr: *mut u8) -> () {
    let c_headers = CString::new(string_content).unwrap();

    let bytes = c_headers.as_bytes_with_nul();

    let header_bytes = std::slice::from_raw_parts_mut(ptr, 1024);

    header_bytes[..bytes.len()].copy_from_slice(bytes);
}

// fn get_return_string(content: String, ptr: *mut u8) {
//     // Assuming ptr is a valid pointer to a buffer with enough space
//     let bytes = content.as_bytes();
//     unsafe {
//         std::ptr::copy_nonoverlapping(bytes.as_ptr(), ptr, bytes.len());
//         *ptr.add(bytes.len()) = 0; // Null-terminate the string
//     }
// }