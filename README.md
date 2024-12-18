## This repository has been archived. Please refer [Rubix Wasm](https://github.com/rubixchain/rubix-wasm) for updated information on Rubix Smart Contracts.


---

**Rubix Smart Contracts**

Smart contracts represent sophisticated business logic encapsulated in machine-readable formats, through programming languages. These contracts are executed within a network’s nodes, operating in a deterministic, sandboxed environment.
At Rubix, we perceive smart contracts not just as mere code but as a specialized form of Programmable Non-Fungible Tokens (NFTs) possessing a dynamic state. Every invocation of a contract function leads to an update in this state, which is recorded and preserved on the Contract Token Chain. This dedicated chain ensures transparent and tamper-proof order of each contract execution.

**Rubix WASM based smart contracts for Wider Adoption**

Rubix protocol is committed to improve adaptability of blockchain technology. Along with its revolutionary proof of pledge protocol aided by zero gas fee transactions , Rubix focuses on making dApp deployment and execution easier for our ecosystem. With WebAssembly(WASM) based smart contracts, existing web2 codebases and developers can migrate their codebase and knowledge into Rubix with ease.
WebAssembly (WASM) is a binary instruction format that allows code to be executed at near-native speed in a safe , sandboxed and deterministic manner across different platforms. Smart contracts can be written in languages that compile to WebAssembly, such as Rust and C/C++, and then executed on a blockchain platform that supports WASM.


**Limitations of using Web Assembly**

1. Linear Memory Model: WebAssembly uses a linear memory model, where the memory is a single, contiguous array of bytes. While this simplicity aids performance, it can be a limitation for certain types of applications. For example, dynamic data structures like linked lists may require additional effort to implement efficiently in a linear memory space.

2. No Direct DOM Access: WebAssembly operates in a sandboxed environment within the web browser and does not have direct access to the Document Object Model (DOM). Interaction with the DOM is typically done through JavaScript, and communication between WebAssembly and JavaScript is necessary for manipulating the DOM.

3. Manual Memory Management: WebAssembly does not have a built-in garbage collector, and memory management is manual. Developers need to allocate and deallocate memory explicitly, which can lead to memory-related bugs if not done carefully.

4. Limited Standard Library: WebAssembly does not have a comprehensive standard library like many higher-level languages. This means that developers often need to rely on external libraries or write more code to implement common functionality that might be readily available in other languages.

5. Debugging Challenges: Debugging WebAssembly code can be more challenging than debugging higher-level languages. While tools like source maps can help map WebAssembly instructions back to the original source code, the debugging experience may not be as seamless as with languages like JavaScript.

6. Bundle Size: WebAssembly binaries can contribute to the overall size of the web application. While efforts have been made to reduce the size of WebAssembly modules, developers need to be mindful of the impact on page load times, especially in contexts where bandwidth is a concern.

7. Limited Multithreading Support: While WebAssembly has some support for parallelism through threads, it currently lacks comprehensive multithreading support. The threading model is cooperative, and true parallelism is not guaranteed. However, proposals for adding more advanced threading features are under consideration.

8. No File System Access: WebAssembly does not have direct access to the file system for security reasons. File system operations must be performed through browser APIs or other mechanisms provided by the host environment.

**How we are doing it ?**

This is our basic or primitive version of the smart contract. As I have mentioned above WebAssembly or wasm has a lot of limitations as well as great potential. So the hardest part here was the memory access part. Wasm has a linear memory and also wasm doesn't have file system access, so these two are the limitations that we needed to tackle.

Wasm does not support any string or complex data types they only supports byte arrays. So what we are doing here is that, all the data which we need to pass to the wasm module is converted into bytes and is appended together to form a large byte array, and the length corresponding to each byte array is also passed along with this to the wasm module. Inside the wasm module, which is written in Rust, we have implemented the logic to decode these byte array to the corresponding data type using the byte array and the length which we have given as input. Then we perform the actions which needs to be done. Here in our case we are using [**Wasmtime**](https://github.com/bytecodealliance/wasmtime-go) as the wasm runtime. 

If you see anything which needs to be changed, feel free to create an issue.
