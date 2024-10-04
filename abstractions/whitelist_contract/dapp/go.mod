module dapp

go 1.21.9

replace wasm_go => ../../wasm_go

require (
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	wasm_go v0.0.0-00010101000000-000000000000
)

require github.com/bytecodealliance/wasmtime-go v1.0.0 // indirect
