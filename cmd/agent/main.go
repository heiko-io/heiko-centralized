package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"

	heiko_rpc "github.com/heiko-io/heiko/internal/rpc"
)

func main() {
	fmt.Println("Hello from the agent!")
	pack, err := ioutil.ReadFile("resources/nix/go.nix")
	if err != nil {
		fmt.Println(err.Error())
	}
	job := heiko_rpc.Job{
		Package: pack,
		Cmd:     "nix-shell go.nix",
		Runtime: "nix",
	}

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

    var result string
    err = client.Call("Job.SubmitJob", job, &result)
    if err != nil {
        log.Fatal("submit job error:", err)
    }
    fmt.Println(result)
}
