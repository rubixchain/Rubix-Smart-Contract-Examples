package contract

import (
	"fmt"
	"io/ioutil"
	"os"

	wasm "github.com/bytecodealliance/wasmtime-go"
	"github.com/joho/godotenv"
)

type ContractExecution struct {
	wasmPath        string
	stateFile       string
	initialised     bool
	pointerPosition int
	instance        *wasm.Instance
	store           *wasm.Store
	memory          *wasm.Memory
}

type Action struct {
	Function string        `json:"function"`
	Args     []interface{} `json:"args"`
}

/*Different functions which we have written here are called from wasm:
  1. "alloc"
  2. "apply_state"
  3. "get_state"
  So there will be a corresponding function with the same name in the rust code too. If the function name
  is different then the system will thorugh an error.
  The initial idea is to make all the mandatory functions into a package which can be easily imported and
  utilised while maintaining the standard for execution.*/

func NewContractExecution(contractId string, port string) (*ContractExecution, error) {
	fmt.Println("Port", port)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	fmt.Println("Contract ID", contractId)
	path := os.Getenv(port) + "SmartContract/"
	c := &ContractExecution{
		wasmPath:  fmt.Sprintf(path+"%s/binaryCodeFile.wasm", contractId),
		stateFile: fmt.Sprintf(path+"%s/schemaCodeFile.json", contractId),
	}
	fmt.Println("Path is ", path)
	fmt.Println("ContractExecution:", c)
	wasmBytes, err := os.ReadFile(c.wasmPath)
	if err != nil {
		fmt.Println("failed to read wasm path,err:", err)
		return nil, err
	}
	c.store = wasm.NewStore(wasm.NewEngine())
	if c.store == nil {
		fmt.Println("not able to create a new store")
	}
	fmt.Println("c.store", c.store.Engine, c.store)
	// wasm1, err := wasm.Wat2Wasm(`
	// 	(module
	//     (import "__wbindgen_placeholder__" "__wbindgen_describe" (func $_ZN12wasm_bindgen19__wbindgen_describe17hdb4530e20704bb46E (type $t0)))
	// 	(import "__wbindgen_placeholder__" "__wbindgen_throw" (func $_ZN12wasm_bindgen16__wbindgen_throw17he65f6939c6dd1424E (type $t12)))
	//     (import "__wbindgen_externref_xform__" "__wbindgen_externref_table_grow" (func $_ZN12wasm_bindgen9externref31__wbindgen_externref_table_grow17hdf1c42fb5204abd7E (type $t1)))
	//     (import "__wbindgen_externref_xform__" "__wbindgen_externref_table_set_null" (func $_ZN12wasm_bindgen9externref35__wbindgen_externref_table_set_null17hd70d97eee4521284E (type $t0)))
	// 	(func (export "run")
	// 	(call $_ZN12wasm_bindgen19__wbindgen_describe17hdb4530e20704bb46E)
	// 	(call $_ZN12wasm_bindgen16__wbindgen_throw17he65f6939c6dd1424E)
	// 	(call $_ZN12wasm_bindgen9externref31__wbindgen_externref_table_grow17hdf1c42fb5204abd7E)
	// 	(call $_ZN12wasm_bindgen9externref35__wbindgen_externref_table_set_null17hd70d97eee4521284E)
	// 	)
	//   `)

	module, err := wasm.NewModule(c.store.Engine, wasmBytes)
	if err != nil {
		fmt.Println("failed to compile new wasm module,err:", err)
		return nil, err
	}
	//item := wasm.WrapFunc(c.store, )

	instance, err := wasm.NewInstance(c.store, module, nil)
	if err != nil {
		fmt.Println("failed to instantiate wasm module,err:", err)
		return nil, err
	}

	allocFn := instance.GetExport(c.store, "alloc").Func()
	address, err := allocFn.Call(c.store)
	if err != nil {
		fmt.Println("failed to alloc wasm mem,err:", err)
		return nil, err
	}

	c.pointerPosition = int(address.(int32))

	c.instance = instance
	c.memory = instance.GetExport(c.store, "memory").Memory()
	c.initialised = true
	fmt.Println("Pointer:", c.pointerPosition)
	//c.apply_state()

	return c, nil
}

func (c *ContractExecution) write(str []byte) int {
	if !c.initialised {
		panic("Contract not initialised")
	}
	ptr := c.pointerPosition

	fmt.Print("Writing to memory: ")
	fmt.Println(str)

	fmt.Print("Pointer position: ")
	fmt.Println(ptr)

	copy(
		c.memory.UnsafeData(c.store)[ptr:],
		str,
	)

	c.pointerPosition += len(str) + 1
	fmt.Println("Latest pointer position", c.pointerPosition)
	return ptr
}

func (c *ContractExecution) readAtCurrentPointer() string {
	if !c.initialised {
		panic("Contract not initialised")
	}

	pointer := c.pointerPosition
	view := c.memory.UnsafeData(c.store)[pointer:]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}

	str := string(view[:length])
	c.pointerPosition += length + 1
	return str
}

func (c *ContractExecution) ReadStateFile() string {
	if !c.initialised {
		panic("Contract not initialised")
	}

	file, err := os.ReadFile(c.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}

		panic(err)
	}

	return string(file)
}

// func (c *ContractExecution) apply_state() {
// 	if !c.initialised {
// 		panic("Contract not initialised")
// 	}

// 	state := c.ReadStateFile()
// 	if state != "" {
// 		pointer := c.write(state)
// 		c.instance.GetExport(c.store, "apply_state").Func().Call(c.store, pointer)
// 	}
// }

func (c *ContractExecution) ProcessActions(actions []Action, jsonStr []byte) {
	if !c.initialised {
		panic("Contract not initialised")
	}

	fmt.Println("The given json string ", jsonStr)
	for _, action := range actions {
		// map on action.args and store to pointers
		pointers := make([]interface{}, len(action.Args))
		// for i, arg := range action.Args {
		// 	pointers[i] = c.write(arg.(string))
		// }
		pointers[0] = c.write(jsonStr)
		functionRef := c.instance.GetExport(c.store, action.Function)
		fmt.Println(functionRef)
		fmt.Println("Function", action.Function)
		functionRef.Func().Call(c.store, pointers...)
	}

	c.save_state()
}

func (c *ContractExecution) save_state() {
	if !c.initialised {
		panic("Contract not initialised")
	}
	fmt.Println("Save State function Called ")
	c.instance.GetExport(c.store, "get_state").Func().Call(c.store, c.pointerPosition)

	state := c.readAtCurrentPointer()
	fmt.Println("State ", state)
	err := ioutil.WriteFile(c.stateFile, []byte(state), 0o644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Save State function Called ")
}
