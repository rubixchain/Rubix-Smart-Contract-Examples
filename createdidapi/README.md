## Bidding-Contract

### Introduction
The Bidding-Contract operates similarly to the previously implemented voting contract. The primary objective is to identify the highest bidder through a streamlined process involving fixed memory allocation and pointer-based operations.

### Process Flow
1. Contract Deployment: The Auctioneer, acting as the initiating node, deploys the contract to initiate the bidding process.
2. Subscription and Participation: Interested nodes subscribe to the contract and place their bids by executing the contract within a specified time frame (e.g., 2 minutes). For demonstration purposes, we will manually trigger the contract once all predefined nodes have executed their bids.
3. Contract Execution: The entire smart contract token chain is passed as a byte array to a WebAssembly (WASM) file, which processes the data to determine the highest bidder.

### Implementation Details
1. Memory Management: The program allocates a fixed memory space and returns the pointer to this memory. All subsequent operations use this pointer, ensuring efficient memory usage.
2. Bid Evaluation: The WASM file processes the byte array containing the bid information and identifies the highest bidder based on the executed contract logic.