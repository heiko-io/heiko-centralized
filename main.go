package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/heiko-io/heiko/internal/agent"
	"github.com/heiko-io/heiko/internal/controller"
	"github.com/heiko-io/heiko/internal/client"
)

var cli struct {
	Submit struct {
		Path string `arg:"" type:"path" optional:""  help:"Path to folder containing code and .heiko.toml" `
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
		client.Submit(cli.Submit.Path)
	case "submit <path>":
		fmt.Println("submiting job...")
		client.Submit(cli.Submit.Path)
	case "agent":
		fmt.Println("starting agent...")
		agent.Start()
	case "controller":
		fmt.Println("starting controller...")
		controller.Start()
	}
}
