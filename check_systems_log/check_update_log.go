package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

func main() {
	server := "<server ip/dns>"
	username := "<username>"
	password := "<password>"

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", server, config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	script := `journalctl | grep -E "apt-get|apt" | grep -E "update|upgrade"`

	output, err := session.CombinedOutput(script)
	if err != nil {
		log.Fatalf("Failed to run script: %s", err)
	}

	fmt.Println("System log from", server, "Linux Server regarding update and upgrade:")
	fmt.Println(string(output))
}
