package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

	datalen int
	data    []byte
}

type Action struct {
	Function string        `json:"function"`
	Args     []interface{} `json:"args"`
}
type rbtTransdata struct {
	Sender     string
	Receiver   string
	TokenCount float64
	Comment    string
	Type       int
	Password   string
	port       string
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

	linker.FuncWrap("env", "rbt_transfer", c.initiateRbtTransfer)
	linker.FuncWrap("env", "create_did", c.initiatecreateDid2)
	// linker.FuncWrap("env", "load_input", c.loadInput)
	// linker.FuncWrap("env", "Addition", c.sum)
	// linker.FuncWrap("env", "Dumpoutput", c.dumpOutput)
	module, err := wasm.NewModule(c.store.Engine, wasmBytes)
	if err != nil {
		fmt.Println("failed to compile new wasm module,err:", err)
		return nil, err
	}

	instance, err := linker.Instantiate(c.store, module)
	// instance, err := wasm.NewInstance(c.store, module, nil)
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
	//c.apply_state()

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

func (c *ContractExecution) readAtCurrentPointer() string {
	if !c.initialised {
		panic("Contract not initialised")
	}

	pointer := c.pointerPosition
	fmt.Println("Pointer position in readAtCurrentPointer function", pointer)
	view := c.memory.UnsafeData(c.store)[pointer:]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length in readAtCurrentPointer function:", length)
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

	c.save_state()
}

func (c *ContractExecution) save_state() {
	if !c.initialised {
		panic("Contract not initialised")
	}
	fmt.Println("Save State function Called ")
	fmt.Println("pointer position in save_state function is:", c.pointerPosition)
	c.instance.GetExport(c.store, "get_state").Func().Call(c.store, c.pointerPosition)
	fmt.Println("pointer position in save_state function after get state function gets called is:", c.pointerPosition)

	state := c.readAtCurrentPointer()
	fmt.Println("State ", state)
	err := ioutil.WriteFile(c.stateFile, []byte(state), 0o644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Save State function Called ")
}
func (c *ContractExecution) initiateRbtTransfer() {
	data := rbtTransdata{
		Sender:     "bafybmihsa7qc5onikjlxvguxifnh7xz7t57q4mqnopee62geheno4iia2m",
		Receiver:   "bafybmibpgv4fe4xr7wwolrymxfphe7o45r4mynnzam6ohqqzvh3usmue2e",
		TokenCount: 1.0,
		Comment:    "Payment for services",
		Type:       2,
		Password:   "mypassword",
		port:       "20012",
	}
	// rbtJSON,err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println("error in marshaling JSON:", err)
	// 	return
	// }
	rbtTransfer(data.Sender, data.Receiver, data.TokenCount, data.Comment, data.Type, data.Password, data.port)
	// c.rbtdata = rbtJSON
	// copy(c.memory.UnsafeData(c.store)[pointer:pointer+int32(len(c.rbtdata))], c.rbtdata)

}
func rbtTransfer(Sender string, Receiver string, TokenCount float64, Comment string, Type int, Password string, port string) {
	data := map[string]interface{}{
		"receiver":   Receiver,
		"sender":     Sender,
		"tokenCOunt": TokenCount,
		"comment":    Comment,
		"type":       Type,
		"password":   Password,
	}
	bodyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error in marshaling JSON:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%s/api/initiate-rbt-transfer", port)
	// "/api/initiate-rbt-transfer"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	data2, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}
	// Process the data as needed
	fmt.Println("Response Body in rbtTransfer :", string(data2))
	var response map[string]interface{}
	err3 := json.Unmarshal(data2, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response:", err3)
	}

	result := response["result"].(map[string]interface{})
	id := result["id"].(string)
	SignatureResponse(id, port)

	defer resp.Body.Close()
}

//	    Type              int    `json:"type"`
//		Dir               string `json:"dir"`
//		Config            string `json:"config"`
//		RootDID           bool   `json:"root_did"`
//		MasterDID         string `json:"master_did"`
//		Secret            string `json:"secret"`
//		PrivPWD           string `json:"priv_pwd"`
//		QuorumPWD         string `json:"quorum_pwd"`
//		ImgFile           string `json:"img_file"`
//		DIDImgFileName    string `json:"did_img_file"`
//		PubImgFile        string `json:"pub_img_file"`
//		PrivImgFile       string `json:"priv_img_file"`
//		PubKeyFile        string `json:"pub_key_file"`
//		PrivKeyFile       string `json:"priv_key_file"`
//		QuorumPubKeyFile  string `json:"quorum_pub_key_file"`
//		QuorumPrivKeyFile string `json:"quorum_priv_key_file"`
//		MnemonicFile      string `json:"mnemonic_file"`
//		ChildPath         int    `json:"childPath"`
//	}

func (c *ContractExecution) initiatecreateDid(pointer int32) {

	data := CreateDidData{
		Type: 0,
		Dir:  "",
		// Config
		// RootDID
		// MasterDID
		// Secret
		PrivPWD: "mypassword",
		// QuorumPWD:      "mypassword",
		ImgFile:        "/home/rubix/Sai-Rubix/rubixgoplatform/linux/image.png",
		DIDImgFileName: "",
		// PubImgFile:     "",
		// PrivImgFile
		// PubKeyFile: "",
		// PrivKeyFile
		// QuorumPubKeyFile
		// QuorumPrivKeyFile
		// MnemonicFile
		// ChildPath
	}
	marshalData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error in marshaling JSON:", err)
		return
	}
	copy(c.memory.UnsafeData(c.store)[pointer:pointer+int32(len(marshalData))], marshalData)
	c.datalen = len(marshalData)
	fmt.Println("data length is:", c.datalen)
	view := c.memory.UnsafeData(c.store)[pointer:]
	length := 0
	for _, byte := range view {
		if byte == 0 {
			break
		}
		length++
	}
	fmt.Println("length in readAtCurrentPointer function:", length)
	str := string(view[:length])
	fmt.Println("data in initiatecreateDid function is:", str)
	// c.rbtdata = rbtJSON
	// copy(c.memory.UnsafeData(c.store)[pointer:pointer+int32(len(c.rbtdata))], c.rbtdata)

}
func (c *ContractExecution) initiatecreateDid2(pointer int32) {
	c.initiatecreateDid(pointer)
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
	var response CreateDidData
	err3 := json.Unmarshal(c.data, &response)
	if err3 != nil {
		fmt.Println("Error unmarshaling response in initiatecreateDid2:", err3)
	}
	port := "20009"
	createDid(response, port)

}

// func createDid(Type int,Dir string,Config string,RootDID bool,MasterDID string,Secret string,PrivPWD string,QuorumPWD string,ImgFile string,DIDImgFileName string,PubImgFile string,PrivImgFile string,PubKeyFile string,PrivKeyFile string,QuorumPubKeyFile string,QuorumPrivKeyFile string,MnemonicFile string,ChildPath int){
// func createDid(Type int, Dir string, PrivPWD string, QuorumPWD string, ImgFile string, DIDImgFileName string, PubImgFile string, PubKeyFile string, port string) {
func createDid(data CreateDidData, port string) {
	// Create a buffer to hold the multipart form data
	var requestBody bytes.Buffer

	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Add form fields
	writer.WriteField("type", fmt.Sprintf("%d", data.Type))
	writer.WriteField("dir", data.Dir)
	writer.WriteField("priv_pwd", data.PrivPWD)
	// writer.WriteField("quorum_pwd", QuorumPWD)
	writer.WriteField("did_img_file", data.DIDImgFileName)
	writer.WriteField("pub_img_file", data.PubImgFile)

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
