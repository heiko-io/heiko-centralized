package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/rpc"
	"os"

	heiko_rpc "github.com/heiko-io/heiko/internal/rpc"
)

func Start() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Hello from the controller!")
	jobs := new(heiko_rpc.RpcController)
	jobs.Queue = make(chan heiko_rpc.Job, 100)
	jobs.Client = client
	go schedule_jobs(jobs)
	listen(jobs)
}

func listen(jobs *heiko_rpc.RpcController) {
	rpc.Register(jobs)
	rpc.HandleHTTP()
	fmt.Println("Listening on port 8080!")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func schedule_jobs(jobs *heiko_rpc.RpcController) {
	for {
		job := <-jobs.Queue
		if find_free_node() {
			var reply string
			jobs.Client.Call("Job.RunJob", job, &reply)
			fmt.Println(reply)
		}
	}
}

// should find if a node is free and return that
// for now it's a hack that returns 1 / 0
func find_free_node() bool {
	return rand.Int() % 2 == 0
}
