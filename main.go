package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lrwx00t/go-remoteexec/ssh_utils"
	"golang.org/x/crypto/ssh"
)

// TODO
// upstream / downstream bootstrap mixed?
// yaml
// template project
// logging for execution
// refactor

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
	}

	remote_server_full_address := fmt.Sprintf("%s:22", *ssh_remote_address_flag)
	conn, err := ssh.Dial("tcp", remote_server_full_address, sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: %s", err)
		os.Exit(1)
	}
	defer conn.Close()
	delete_repo := "rm -fR go-bootstrap"
	clone_repo := "git clone https://github.com/0xack13/go-bootstrap"
	session_execute(delete_repo, conn)
	session_execute(clone_repo, conn)
	session_execute("cd go-bootstrap && make install", conn)
	fmt.Println("finished boostrap. cleaning up..")
	session_execute(delete_repo, conn)

}

func session_execute(cmd string, conn *ssh.Client) {
	session, err := conn.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %s", err)
		os.Exit(1)
	}
	defer session.Close()

	out, err := session.Output(cmd)
	if err != nil {
		fmt.Printf("Failed to run command: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%s", out)
}
