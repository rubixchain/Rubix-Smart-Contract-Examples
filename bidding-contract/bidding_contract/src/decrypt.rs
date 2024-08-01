use ecies::{decrypt, encrypt, utils::generate_keypair};
use std::io::{self, Read};
use std::env;
use std::fs;
use pem::parse;
use aes_gcm::aead::{Aead, generic_array::GenericArray};
use aes_gcm::Aes256Gcm;
use sha2::{Sha256, Digest};
use std::error::Error;
use std::fmt;
//use crate::seal;
pub fn decrypt_smartcontract_data(decrypted_key:&[u8],ciphertext:Vec<u8>)->Vec<u8>{
    // let read_enc_privkey = fs::read_to_string(enc_privkey_path).expect("not able to read the private key");
    // println!("private key read from the file is {}",read_enc_privkey);
    // let read_enc_privkey_bytes:Vec<u8> = read_enc_privkey.into();
    // let key = "mypassword";
    // let pem_enc_privkey = parse(read_enc_privkey_bytes).expect("not able to read from the pem file");
    // let enc_privkey_decoded_bytes_vec = pem_enc_privkey.contents;
    // let enc_privkey_bytes: &[u8] = &enc_privkey_decoded_bytes_vec;

    #[cfg(not(feature = "x25519"))]
    //let decrypted_key = seal::unseal(key, &enc_privkey_decoded_bytes_vec).expect("not able to unseal the encrypted key"); 
    //println!("Unsealed pvtkey is {:?}",decrypted_key);
    //let test_enc_privkey :[u8;32]= [26, 74, 206, 110, 150, 148, 87, 32, 213, 102, 150, 120, 224, 105, 131, 103, 58, 95, 72, 72, 142, 240, 97, 25, 113, 39, 140, 138, 164, 82, 187, 147];
    //let decrypted_key = [26, 74, 206, 110, 150, 148, 87, 32, 213, 102, 150, 120, 224, 105, 131, 103, 58, 95, 72, 72, 142, 240, 97, 25, 113, 39, 140, 138, 164, 82, 187, 147];
    let decrypted_message = decrypt(decrypted_key, &ciphertext).expect("not able to decrypt");
    //let decrypted_msg_string = String::from_utf8(decrypted_message).unwrap();
    
    println!("decrypted message is {:?}",decrypted_message);
    return decrypted_message
    
    
    }
