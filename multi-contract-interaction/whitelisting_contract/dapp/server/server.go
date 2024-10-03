package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const API_RUN_WHITELISTING_CONTRACT = "/run-whitelisting-contract"
const API_GET_WHITELISTED_DID = "/get-whitelisted-dids"
const API_SYNC_STATE = "/sync-state"

const WHITELISTING_CONTRACT = ""

func RunServer() {
	fmt.Println("DApp Server Started")
	r := mux.NewRouter()

	r.HandleFunc(API_RUN_WHITELISTING_CONTRACT, runWhitelistContractHandle).Methods("POST")
	r.HandleFunc(API_GET_WHITELISTED_DID, getWhitelistedDidsHandler).Methods("GET")
	r.HandleFunc(API_SYNC_STATE, syncStateHandler).Methods("POST")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}