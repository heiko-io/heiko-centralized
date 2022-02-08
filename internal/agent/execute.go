package agent

import (
	"fmt"
	"log"
	"os/exec"

	heiko_rpc "github.com/heiko-io/heiko/internal/rpc"
)

func ExecuteJobs(execJobChan chan heiko_rpc.Job) {
	for {
		job := <-execJobChan
		go runJob(job)
	}
}

func runJob(job heiko_rpc.Job) {
	cmd := exec.Command("nix-shell", "shell.nix", "--run", `"./run.sh"`)
	cmd.Dir = job.Name
	// debug
	fmt.Println("Running command: ", cmd.String())
	// test_cmd := exec.Command("pwd")
	// test_cmd.Dir = job.Name
	// test_out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Fatalln("Failed to get pwd: ", err)
	// }
	// fmt.Println("pwd: ", test_out)


	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("Failed to run job: ", err, ": ", string(out))
	}

	fmt.Println("Result: ", string(out))
}
