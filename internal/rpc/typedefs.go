package rpc

import (
	"net/rpc"
	"github.com/heiko-io/heiko/resources/nix"
)


var Runtime = nix.GetRuntimes()

type Job struct {
	Runtime string
	Package []byte
	Cmd     []string
	Init    []string
	Name    string
}

type RpcController struct {
	Queue  chan Job
	Client *rpc.Client
}
