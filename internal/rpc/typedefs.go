package rpc

type Job struct {
	Runtime string
	Package []byte
	Cmd     string
}
