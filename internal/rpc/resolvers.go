package rpc

import "fmt"

func (t *Job) RunJob(job *Job, reply *string) error {
	*reply = fmt.Sprintf("Running: Job { Runtime: %s, Cmd: %s }", job.Cmd, job.Runtime)
	return nil
}

func (t *RpcController) SubmitJob(job *Job, reply *string) error {
	t.Queue <- *job
	return nil
}
