package main

import (
	"dapp/server"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Server has been started")
	go server.RunServer()

	time.Sleep(300000 * time.Second)
	fmt.Println("Server has stopped!")
}
