#!/bin/sh

cargo build --target wasm32-unknown-unknown && cp target/wasm32-unknown-unknown/debug/addition_contract.wasm ../addition_app/