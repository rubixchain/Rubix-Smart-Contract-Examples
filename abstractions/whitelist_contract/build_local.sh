#!/bin/sh

cargo build --target wasm32-unknown-unknown && cp target/wasm32-unknown-unknown/debug/whitelist_contract.wasm ./artifacts/