package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/heiko-io/heiko/internal/agent"
	"github.com/heiko-io/heiko/internal/controller"
)

var cli struct {
	Submit struct {
		Path string `arg:"" optional:"" help:"Path to folder containing code and .heiko.toml" type:"path"`
	} `cmd:"" help:"Submit a job to heiko"`
	Agent      struct{} `cmd:"" help:"Spawn a heiko agent"`
	Controller struct{} `cmd:"" help:"Spawn a heiko controller"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("heiko"),
		kong.Description("Heiko is your one stop shop for creating a private serverless cloud infra with minimal effort."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	switch ctx.Command() {
	case "submit":
		fmt.Println("submiting job...")
	case "agent":
		fmt.Println("starting agent...")
		agent.Start()
	case "controller":
		fmt.Println("starting controller...")
		controller.Start()
	}
}
