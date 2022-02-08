package client

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	heiko_rpc "github.com/heiko-io/heiko/internal/rpc"
)

type Config struct {
	Name    string   `toml:"name"`
	Runtime string   `toml:"runtime"`
	Cmds    []string `toml:"cmds"`
	Init    []string `toml:"init"`
}

func ReadConfig(dir string) heiko_rpc.Job {
	var conf Config
	// contents, err := ioutil.ReadFile(dir+"./heiko.toml")
	_, err := toml.DecodeFile(dir+"/.heiko.toml", &conf)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	// fmt.Println(conf.Runtime)
	// fmt.Println(conf)
	if _, ok := heiko_rpc.Runtime[conf.Runtime]; !ok {
		log.Fatalln("Invalid runtime! Choose one of node, python, python3, go or rust")
		os.Exit(1)
	}
	
	return heiko_rpc.Job{
		Runtime: conf.Runtime,
		Name: conf.Name,
		Cmd: conf.Cmds,
		Init: conf.Init,
	}
}
