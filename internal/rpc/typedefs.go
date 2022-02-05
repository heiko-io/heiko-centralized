package rpc

var Runtime = map[string]string{
	"node": "resources/nix/node.nix",
	"python2": "resources/nix/python.nix",
	"python3": "resources/nix/python3.nix",
	"go": "resources/nix/go.nix",
	"rust": "resources/nix/rust.nix",
}

type Job struct {
	Runtime string
	Package []byte
	Cmd     []string
	Init    []string
	Name    string
}
