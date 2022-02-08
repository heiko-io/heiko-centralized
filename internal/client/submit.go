package client

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
)

func Submit(path string) {
	if path == "" {
		fmt.Println("Using current directory......")
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory: %v", err)
			os.Exit(1)
		}
		path = dir
	}
	fmt.Println(path)
	job := ReadConfig(path)
	job.Package = create_tar(path).Bytes()
	client, err := rpc.DialHTTP("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln("Error connecting to host:", err)
		os.Exit(1)
	}
	err = client.Call("RpcController.SubmitJob", job, nil)
	if err != nil {
		log.Fatalln("Error submitting job: ", err)
	}
}

func create_tar(path string) *bytes.Buffer{
	var buf bytes.Buffer
	err := compress(path, &buf)
	if err != nil {
		log.Fatalf("Error compressing directory: %v", err)
	}
	return &buf
}

func compress(src string, buf io.Writer) error {
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	filepath.Walk(src, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Could not walk file %s: %v", path, err)
			return err
		}
	
		og_path := path
		if strings.HasPrefix(path, src) {
			og_path = path[len(src):]
		}
		if len(og_path) == 0 {
			return nil
		}
		if og_path[0] == '/' {
			og_path = og_path[1:]
		}
		fmt.Println(og_path)

		if fi.IsDir() {
			og_path += "/"
		}

		// if the file is a symlink, this will be the linked path
		link_path := og_path
		isSymLink := fi.Mode() & os.ModeSymlink == os.ModeSymlink
		if isSymLink {
			link_path, err = os.Readlink(path)
			if err != nil {
				log.Fatalf("Could not get link path (%s): %v", path, err)
			}
		}


		// generate and write header
		header, err := tar.FileInfoHeader(fi, link_path)
		if err != nil {
			return err
		}

		header.Name = og_path

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// if not a dir write file content
		if !fi.IsDir() && !isSymLink {
			data, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}

	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}

	return nil
}
