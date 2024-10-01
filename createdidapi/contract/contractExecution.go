package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	// "io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	setup "github.com/rubixchain/rubixgoplatform/setup"

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

	data []byte
}

type Action struct {
	Function string        `json:"function"`
	Args     []interface{} `json:"args"`
}

type CreateDidData1 struct {
	Type    int
	PrivPWD string
	ImgFile string
}
type CreateDidData struct {
	Type              int    `json:"type"`
	Dir               string `json:"dir"`
	Config            string `json:"config"`
	RootDID           bool   `json:"root_did"`
	MasterDID         string `json:"master_did"`
	Secret            string `json:"secret"`
	PrivPWD           string `json:"priv_pwd"`
	QuorumPWD         string `json:"quorum_pwd"`
	ImgFile           string `json:"img_file"`
	DIDImgFileName    string `json:"did_img_file"`
	PubImgFile        string `json:"pub_img_file"`
	PrivImgFile       string `json:"priv_img_file"`
	PubKeyFile        string `json:"pub_key_file"`
	PrivKeyFile       string `json:"priv_key_file"`
	QuorumPubKeyFile  string `json:"quorum_pub_key_file"`
	QuorumPrivKeyFile string `json:"quorum_priv_key_file"`
	MnemonicFile      string `json:"mnemonic_file"`
	ChildPath         int    `json:"childPath"`
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
		stateFile: fmt.Sprintf(path+"%s/SchemaCodeFile.json", contractId),
	}
	fmt.Println("Path is ", path)
	fmt.Println("ContractExecution:", c)
	wasmBytes, err := os.ReadFile(c.wasmPath)
	if err != nil {
		return nil, err
	}

	engine := wasm.NewEngine()
	linker := wasm.NewLinker(engine)
	linker.DefineWasi()

	c.store = wasm.NewStore(engine)
	if c.store == nil {
		fmt.Println("not able to create a new store")
	}
	fmt.Println("c.store", c.store.Engine, c.store)

	linker.FuncWrap("env", "create_did", c.initiatecreateDid2)

	module, err := wasm.NewModule(c.store.Engine, wasmBytes)
	if err != nil {
		fmt.Println("failed to compile new wasm module,err:", err)
		return nil, err
	}
	instance, err := linker.Instantiate(c.store, module)
	if err != nil {
		fmt.Println("failed to instantiate wasm module,err:", err)
		return nil, err
	}
	allocFn := instance.GetExport(c.store, "alloc").Func()
	address, err := allocFn.Call(c.store)
	if err != nil {
		return nil, err
	}

	c.pointerPosition = int(address.(int32))

	c.instance = instance
	c.memory = instance.GetExport(c.store, "memory").Memory()
	c.initialised = true
	fmt.Println("Pointer:", c.pointerPosition)

	return c, nil
}

func (c *ContractExecution) write(str string) int {
	if !c.initialised {
		panic("Contract not initialised")
	}
	ptr := c.pointerPosition
	fmt.Print("length of the string is:", len(str))
	fmt.Print("\n Writing to memory: ")
	fmt.Println(str)

	fmt.Print("Pointer position: ")
	fmt.Println(ptr)

	copy(
		c.memory.UnsafeData(c.store)[ptr:],
		[]byte(str),
	)

	c.pointerPosition += len(str) + 1
	fmt.Println("Latest pointer position", c.pointerPosition)
	return ptr
}

func (c *ContractExecution) ProcessActions(actions []Action, jsonStr string) {
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
		fmt.Println("Pointers in ProcessActions function is:", pointers)
		functionRef := c.instance.GetExport(c.store, action.Function)
		fmt.Println(functionRef)
		fmt.Println("Function", action.Function)
		functionRef.Func().Call(c.store, pointers...)
	}

	// c.save_state()
}
func (c *ContractExecution) initiatecreateDid2(pointer int32) {
	// c.initiatecreateDid(pointer)
	// copy(c.data, c.memory.UnsafeData(c.store)[pointer:pointer+int32(c.datalen)])
	view := c.memory.UnsafeData(c.store)[pointer:]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length in initiatecreateDid2 func is :", length)
	str := string(view[:length])
	c.data = view[:length]
	fmt.Println("data in initiatecreateDid2 function is:", str)
	fmt.Println("data in initiatecreateDid2 is: ", string(c.data))
	var response1 CreateDidData1
	err3 := json.Unmarshal(c.data, &response1)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in initiatecreateDid2:", err3)
	}
	fmt.Println("Unmarshalled data in initiatecreateDid2 func is:", response1)
	var response CreateDidData
	response.Type = response1.Type
	response.PrivPWD = response1.PrivPWD
	response.ImgFile = response1.ImgFile
	port := "20009"
	createDid(response, port)

}

func createDid(data CreateDidData, port string) {
	// Create a buffer to hold the multipart form data
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)
	fmt.Println("Printing the data in CreateDid function", data)
	// Add form fields
	writer.WriteField("type", fmt.Sprintf("%d", data.Type))
	// writer.WriteField("dir", data.Dir)
	writer.WriteField("priv_pwd", data.PrivPWD)
	// writer.WriteField("quorum_pwd", QuorumPWD)
	// writer.WriteField("did_img_file", data.DIDImgFileName)
	// writer.WriteField("pub_img_file", data.PubImgFile)
	writer.WriteField("image_file", data.ImgFile)
	// Add the image file to the form
	fmt.Println("Image file name is:", data.ImgFile)
	file, err := os.Open(data.ImgFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	jd, err := json.Marshal(&data)
	if err != nil {
		fmt.Println("Failed to parse json data", "err", err)
		// return "Failed to parse json data", false
	}
	fields := make(map[string]string)
	fields[setup.DIDConfigField] = string(jd)
	writer.WriteField("did_config", string(jd))
	// files := make(map[string]string)
	formFile, err := writer.CreateFormFile("img_file", data.ImgFile)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		fmt.Println("Error copying file content:", err)
		return
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing multipart writer:", err)
		return
	}

	// Create the request URL
	url := fmt.Sprintf("http://localhost:%s/api/createdid", port)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Set the Content-Type header to multipart/form-data with the correct boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Read and print the response body
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(data2))
}
