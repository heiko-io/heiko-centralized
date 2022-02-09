package rpc

import "fmt"

// TODO: can be renamed delegate job or something? Since it doesn't actually run it
func (t *RpcController) RunJob(job *Job, reply *string) error {
	*reply = fmt.Sprintf("Running: Job { Name: %s, Runtime: %s, Cmd: %s }",
		job.Name, job.Cmd, job.Runtime)
	t.Queue <- *job
	return nil
}

func (t *RpcController) SubmitJob(job *Job, reply *string) error {
	t.Queue <- *job
	return nil
}
