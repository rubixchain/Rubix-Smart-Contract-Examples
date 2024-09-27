use proc_macro::TokenStream;
use quote::quote;
use syn::{parse_macro_input, ItemFn};

#[proc_macro_attribute]
pub fn wasm_export(_attr: TokenStream, item: TokenStream) -> TokenStream {
    // Parse the input function
    let input = parse_macro_input!(item as ItemFn);

    // Extract the function name and input/output types
    let func_name = &input.sig.ident;
    let input_type = match input.sig.inputs.first() {
        Some(syn::FnArg::Typed(arg)) => &arg.ty,
        _ => panic!("Expected a function with one argument"),
    };
    let output_type = match &input.sig.output {
        syn::ReturnType::Type(_, ty) => ty,
        _ => panic!("Expected a function with a return type"),
    };

    // Generate a new name for the wrapper function
    let wrapper_func_name = syn::Ident::new(&format!("{}_wrapper", func_name), func_name.span());

    // Generate the wrapper function
    let expanded = quote! {
        // Original function
        #input

        // Generated wrapper function
        #[no_mangle]
        pub extern "C" fn #wrapper_func_name(input_ptr: *mut u8, input_len: usize, output_ptr_ptr: *mut *mut u8, output_len_ptr: *mut usize) -> i32 {
            use std::slice;
            use std::ptr;
            use serde::{Serialize, Deserialize};
            use serde_json;

            // Deserialize input data
            let input_data = unsafe { slice::from_raw_parts(input_ptr, input_len) };
            let input: #input_type = match serde_json::from_slice(input_data) {
                Ok(data) => data,
                Err(_) => return 1,
            };

            // Call the original function
            let result: #output_type = #func_name(input);

            // Serialize output data
            let serialized_output = match serde_json::to_vec(&result) {
                Ok(data) => data,
                Err(_) => return 1,
            };

            // Allocate memory for output data
            let output_len = serialized_output.len();
            let output_ptr = alloc(output_len);

            // Write serialized data to output_ptr
            unsafe {
                ptr::copy_nonoverlapping(serialized_output.as_ptr(), output_ptr, output_len);
                *output_ptr_ptr = output_ptr;
                *output_len_ptr = output_len;
            }

            0
        }
    };

    TokenStream::from(expanded)
}
