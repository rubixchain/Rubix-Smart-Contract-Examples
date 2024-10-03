# Multi Contract Interaction

In this example, we will draw inspiration from the Bidding Contract and build an example to demonstrate multi-contract interaction. Please note that the implementation would be simplistic as the purpose of this demostration is to showcase the interaction of 

## Contract A - Whitelisting 

It will be a whitelisting contract where we will we store a list of DIDs who will be eligible to place the Bid for the bidding contract. 

## Contract B - Bidding Contract

The bidding contract will take the Bidding info as follows:

```json
{
    "did": "did123..",
    "amount": "40.0"
}
```

It will then check if the input `did` is present in Whitelisting contract. If yes, then the `did` will be allowed to place the Bid in the bidding, else it will be skipped 

