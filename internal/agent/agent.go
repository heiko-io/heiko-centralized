package agent

import (
	"fmt"
	"net/http"
	"net/rpc"

	heiko_rpc "github.com/heiko-io/heiko/internal/rpc"
)

func Start() {
	fmt.Println("Hello from the agent!")
	job := new(heiko_rpc.Job)
	rpc.Register(job)
	rpc.HandleHTTP()
	fmt.Println("Listening on port 1234!")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
