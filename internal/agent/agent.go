package agent

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	heiko_rpc "github.com/heiko-io/heiko/internal/rpc"
)

func Start() {
	fmt.Println("Hello from the agent!")
	jobs := new(heiko_rpc.RpcController)
	jobs.Queue = make(chan heiko_rpc.Job, 100)
	rpc.Register(jobs)
	rpc.HandleHTTP()
	fmt.Println("Listening on port 1234!")

	execJobChan := make(chan heiko_rpc.Job, 100)

	go runJobs(jobs, execJobChan)

	go ExecuteJobs(execJobChan)

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// run jobs as they come
func runJobs(jobs *heiko_rpc.RpcController, execJobChan chan heiko_rpc.Job) {
	for {
		job := <-jobs.Queue
		path := savePackage(job.Package, job.Name, job.Runtime)
		createRunScript(job, path)
		job.Name = path
		execJobChan <- job
	}
}

// save the tar package received as bytes
func savePackage(tar_file []byte, name string, runtime string) string {
	path := fmt.Sprintf("/tmp/%s-%d", name, rand.Int())
	fmt.Println("Saving package to path: ", path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to created dir %s: %v", path, err)
	}

	err = decompress(bytes.NewBuffer(tar_file), path)
	if err != nil {
		log.Fatalf("Failed to write package to %s: %v", path, err)
	}

	saveNixFile(runtime,path)

	return path
}

// get the correct nix environment file to use and run
func saveNixFile(runtime string, path string) {
	nix_file := heiko_rpc.Runtime[runtime]
	file_name := fmt.Sprintf("%s/shell.nix", path)
	err := ioutil.WriteFile(file_name, nix_file, os.ModePerm)
	if err != nil {
		log.Fatalf("Error writing nix file to path %s: %v", path, err)
	}
}

// create a run.sh file that would be used to build and/or execute the job
func createRunScript(job heiko_rpc.Job, path string) {
	file, err := os.Create(path+"/run.sh")
	if err != nil {
		log.Fatalf("Failed to create file %v: %v\n", path+"/run.sh", err)
		os.Exit(1)
	}

	defer file.Close()

	// TODO: the same call to WriteString and err handling is repeated many times here
	_, err = file.WriteString("#!/bin/bash\n")
	if err != nil {
		log.Fatalln("Failed to write to run file: ", err)
		os.Exit(1)
	}
	
	// generate build phase commands
	if len(job.Init) > 0 {
		_, err = file.WriteString(`echo "BUILDING....."` + "\n")
		if err != nil {
			log.Fatalln("Failed to write to run file: ", err)
			os.Exit(1)
		}
		for _, cmd := range job.Init {
			_, err = file.WriteString(cmd + "\n")
			if err != nil {
				log.Fatalln("Failed to write to run file: ", err)
				os.Exit(1)
			}
		}
	}

	// generate run phase commands
	if len(job.Cmd) == 0 {
		log.Fatalln("No commands to run!")
		os.Exit(1)
	}
	_, err = file.WriteString(`echo "RUNNING....."` + "\n")
	if err != nil {
		log.Fatalln("Failed to write to run file: ", err)
		os.Exit(1)
	}
	for _, cmd := range job.Cmd {
		_, err = file.WriteString(cmd + "\n")
		if err != nil {
			log.Fatalln("Failed to write to run file: ", err)
			os.Exit(1)
		}
	}

	file.Sync()
	
	// TODO: maybe there's a better way to set exec perms?
	cmd := exec.Command("chmod", "+x", "run.sh")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		// TODO: fix error message
		log.Fatalln("Failed to run job: ", err, ": ", string(out))
	}
}

// check for path traversal and correct forward slashes
func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}

// uncompress the package and save it at the correct location
func decompress(src io.Reader, dst string) error {
	// ungzip
	zr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	// untar
	tr := tar.NewReader(zr)

	// uncompress each element
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		target := header.Name

		// validate name against path traversal
		if !validRelPath(header.Name) {
			return fmt.Errorf("tar contained invalid name error %q", target)
		}

		// add dst + re-format slashes according to system
		target = filepath.Join(dst, header.Name)
		// if no join is needed, replace with ToSlash:
		// target = filepath.ToSlash(header.Name)

		// check the type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it (with 0755 permission)
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// if it's a file create it (with same permission)
		case tar.TypeReg:
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return err
			}
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			fileToWrite.Close()
		}
	}

	//
	return nil
}
