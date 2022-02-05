package rpc

import "fmt"

func (t *Job) SubmitJob(job *Job, reply *string) error {
    *reply = fmt.Sprintf("Job{ Runtime: %s, Cmd: %s }", job.Cmd, job.Runtime)
	return nil
}
