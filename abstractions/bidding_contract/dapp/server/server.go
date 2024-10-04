package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const API_RUN_BIDDING_CONTRACT = "/run-bidding-contract"
const API_GET_BIDS = "/get-bids"

func RunServer() {
	fmt.Println("DApp Server Started")
	r := mux.NewRouter()

	r.HandleFunc(API_RUN_BIDDING_CONTRACT, runBiddingContractHandle).Methods("POST")
	r.HandleFunc(API_GET_BIDS, getBidsHandler).Methods("GET")

	err := http.ListenAndServe(":8081", r)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
