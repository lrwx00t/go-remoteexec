package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/lrwx00t/go-remoteexec/ssh_utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// TODO
// upstream / downstream bootstrap mixed?
// yaml
// template project
// logging for execution
// refactor
// pass user
// file transfer / scp / evaluate secret

func main() {

	var ssh_remote_address_flag = flag.String("a", "", "identity")
	var ssh_identity_flag = flag.String("i", "", "identity")

	flag.Parse()
	remoteexec_pass := os.Getenv("REMOTEEXEC_PASS")
	if len(remoteexec_pass) == 0 && ssh_identity_flag == nil {
		panic("you need to provided a password or private key path")
	}

	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password(remoteexec_pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if len(*ssh_identity_flag) > 0 {
		if strings.HasPrefix(*ssh_identity_flag, "~/") {
			dirname, _ := os.UserHomeDir()
			*ssh_identity_flag = filepath.Join(dirname, (*ssh_identity_flag)[2:])
		}
		pemBytes, err := os.ReadFile(*ssh_identity_flag)
		if err != nil {
			panic(err)
		}

		signer, err := ssh_utils.SignerFromPem(pemBytes, []byte(""))
		if err != nil {
			panic(err)
		}
		sshConfig = &ssh.ClientConfig{
			User:            "root",
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
		}
	} else if len(*ssh_identity_flag) == 0 && len(remoteexec_pass) == 0 {
		fmt.Println("using SSH AUTH SOCK")
		agentSock := os.Getenv("SSH_AUTH_SOCK")
		if agentSock == "" {
			panic("SSH_AUTH_SOCK is not set")
		}
		agentSigners := func() ([]ssh.Signer, error) {
			conn, err := net.Dial("unix", agentSock)
			if err != nil {
				return nil, err
			}
			agentClient := agent.NewClient(conn)
			return agentClient.Signers()
		}

		hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}

		sshConfig = &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeysCallback(agentSigners),
			},
			HostKeyCallback: hostKeyCallback,
		}
	}

	remote_server_full_address := fmt.Sprintf("%s:22", *ssh_remote_address_flag)
	conn, err := ssh.Dial("tcp", remote_server_full_address, sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: %s", err)
		os.Exit(1)
	}
	defer conn.Close()
	// Set environment variables for the session
	session, err := conn.NewSession()
	// session, err := conn.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %s", err)
		os.Exit(1)
	}
	request := &struct {
		Name  string
		Value string
	}{
		Name:  "HELLO",
		Value: "WOWOW HELLO1",
	}

	_, err = session.SendRequest("export", true, ssh.Marshal(request))
	if err != nil {
		panic(err)
	}

	ssh_utils.Session_execute("echo $HELLO", conn)
	// ssh_utils.Session_execute("export HELLO=\"mew2\"", conn)
	ssh_utils.Session_execute("echo $HELLO", conn)
	// ---------------
	ssh_utils.SSHCopyFile(conn, "main.go", "main.go.new")
	delete_repo := "rm -fR go-bootstrap"
	clone_repo := "git clone https://github.com/0xack13/go-bootstrap"
	ssh_utils.Session_execute(delete_repo, conn)
	ssh_utils.Session_execute(clone_repo, conn)
	ssh_utils.Session_execute("cd go-bootstrap && make install", conn)
	fmt.Println("finished boostrap. cleaning up..")
	ssh_utils.Session_execute(delete_repo, conn)

}
