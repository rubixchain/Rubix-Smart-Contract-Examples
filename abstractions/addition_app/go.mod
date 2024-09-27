module addition_app

go 1.21.9

replace wasm_go => ../wasm_go

require wasm_go v0.0.0-00010101000000-000000000000

require (
	github.com/bytecodealliance/wasmtime-go v1.0.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
)
