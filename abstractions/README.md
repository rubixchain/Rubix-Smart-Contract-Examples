# Rubix Smart Contract Libraries (In Development)

A list of libraries to support development of Smart Contract on Rubix

## Example : Multi-Contract Interaction

In this scenario, we have two contracts: [Bidding contract](./bidding_contract/) and [Whitelisting contract](./whitelist_contract/). Please note that this is currently a simplictic and loosely tested example to showcase the ability for multi-contract interaction

`Bidding` contract - Users can essentially place BID with their DIDs and specify the amount that they want to bid. Example contract input representation: `{"place_bid": {"did": "did123", "bid": 200}}`.

`Whitelist` contract - Users can essentially place BID with their DIDs and specify the amount that they want to bid. Example contract input representation: `{"place_bid": {"did": "did123", "bid": 200}}`

### Whitelist Contract

The Whitelist Contract DApp server has three endpoints defined:

```
POST: /run-whitelisting-contract  | Callback URL registered on Rubix which is triggered upon Smart Contract Execute API call

GET: /get-whitelisted-dids | It returns the content of its state JSON file (also utilised in Bidding Contract)

POST: /state-sync | It fetches the state transitions of the Whitelist Contract and reconstructs the state JSON file
```

### Bidding Contract

The Bidding Contract DApp server has one endpoint defined:

```
POST: /run-bidding-contract  | Callback URL registered on Rubix which is triggered upon Smart Contract Execute API call

POST: /get-bids | It returns the content of its state JSON file (also utilised in Bidding Contract)

```

In `bidding_contract/dapp/server/handler.go`, we are calling the `checkIfDIDIsWhitelisted` function inside the `runBiddingContractHandle` handler function, which is basically calling the `/get-whitelisted-dids` API to get the list of whitelisted DIDs. We then check if our input DID (to Bidding contract) is present in on the list or not? If yes, then we proceed with the WASM file execution, else we throw an error from DApp's end.