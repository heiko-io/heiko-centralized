package nix

import (
	"embed"
	"log"
)

//go:embed *.nix
var envs embed.FS


func GetRuntimes() map[string][]byte {
	
	node, err := envs.ReadFile("node.nix")
	if err != nil {
		log.Fatalf("Error reading node.nix: %v", err)
	}
	python2, err := envs.ReadFile("python.nix")
	if err != nil {
		log.Fatalf("Error reading python.nix: %v", err)
	}
	python3, err := envs.ReadFile("python3.nix")
	if err != nil {
		log.Fatalf("Error reading python3.nix: %v", err)
	}
	golang, err := envs.ReadFile("golang.nix")
	if err != nil {
		log.Fatalf("Error reading golang.nix: %v", err)
	}
	rust, err := envs.ReadFile("rust.nix")
	if err != nil {
		log.Fatalf("Error reading rust.nix: %v", err)
	}

	return map[string][]byte {
		"node": node,
		"python2": python2,
		"python3": python3,
		"go": golang,
		"rust": rust,
	}
}
