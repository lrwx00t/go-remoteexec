package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func main() {
	// Load the private key from disk
	// keyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
	// keyBytes, err := ioutil.ReadFile(keyPath)
	// if err != nil {
	// 	panic(err)
	// }
	// key, err := ssh.ParsePrivateKey(keyBytes)
	// if err != nil {
	// 	panic(err)
	// }

	// Set the SSH_AUTH_SOCK environment variable
	agentSock := os.Getenv("SSH_AUTH_SOCK")
	if agentSock == "" {
		panic("SSH_AUTH_SOCK is not set")
	}
	// os.Setenv("SSH_AUTH_SOCK", agentSock)

	// Create a function that wraps agent.NewClient and returns a list of signers
	agentSigners := func() ([]ssh.Signer, error) {
		conn, err := net.Dial("unix", agentSock)
		if err != nil {
			return nil, err
		}
		// defer conn.Close()
		agentClient := agent.NewClient(conn)
		// fmt.Println(agentClient.List())
		return agentClient.Signers()
	}

	// Create a HostKeyCallback function that returns nil (accept any host key)
	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	// Connect to the remote server
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentSigners),
		},
		HostKeyCallback: hostKeyCallback,
	}
	conn, err := ssh.Dial("tcp", "159.89.50.2:22", config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Execute a command on the remote server
	session, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	out, err := session.Output("ls -alth &&  hostname -I")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
