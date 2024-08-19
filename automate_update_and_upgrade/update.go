package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

const (
	sshUser          = "<user>"
	sshPassword      = "<ssh-password>"
	sshHost          = "<server-ip>"
	sshPort          = "22"
	localScriptPath  = "updateupgrade.sh"
	remoteScriptPath = "/tmp/updateupgrade.sh"
)

func main() {
	scpCmd := exec.Command("scp", "-P", sshPort, localScriptPath, fmt.Sprintf("%s@%s:%s", sshUser, sshHost, remoteScriptPath))

	if err := scpCmd.Run(); err != nil {
		log.Fatalf("Error while uploading SCP file: %v", err)
	}

	log.Println("File uploaded successfully!")

	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", sshHost, sshPort), sshConfig)
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	log.Println("Connection to server established successfully!")

	chmodSession, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create SSH session: %v", err)
	}
	defer chmodSession.Close()

	if err := chmodSession.Run(fmt.Sprintf("chmod +x %s", remoteScriptPath)); err != nil {
		log.Fatalf("Failed to change file permissions: %v", err)
	}
	log.Println("File permissions changed successfully!")

	runSession, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create SSH session: %v", err)
	}
	defer runSession.Close()

	var stdout, stderr bytes.Buffer
	runSession.Stdout = &stdout
	runSession.Stderr = &stderr

	log.Println("Running script on server...")
	if err := runSession.Run(fmt.Sprintf("/bin/sh %s", remoteScriptPath)); err != nil {
		log.Printf("Error running script: %v", err)
		log.Printf("Output (stderr): %s", stderr.String())
		log.Fatalf("Output (stdout): %s", stdout.String())
	}

	log.Printf("Script output: %s", stdout.String())
}
