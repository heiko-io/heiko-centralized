package client

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func Submit(path string) {
	if path == "" {
		fmt.Println("Using current directory......")
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory: %v", err)
			os.Exit(1)
		}
		path = dir
	}
	fmt.Println(path)
	job := ReadConfig(path)
	client, err := rpc.DialHTTP("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln("Error connecting to host:", err)
		os.Exit(1)
	}
	client.Call("RpcController.SubmitJob", job, nil)
}

