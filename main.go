package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		signer, err := signerFromPem(pemBytes, []byte(""))
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

func signerFromPem(pemBytes []byte, password []byte) (ssh.Signer, error) {
	err := errors.New("pem decode failed, no key found")
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, err
	}

	if x509.IsEncryptedPEMBlock(pemBlock) {
		pemBlock.Bytes, err = x509.DecryptPEMBlock(pemBlock, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("decrypting PEM block failed %v", err)
		}

		key, err := parsePemBlock(pemBlock)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("creating signer from encrypted key failed %v", err)
		}

		return signer, nil
	} else {
		// generate signer instance from plain key
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return nil, fmt.Errorf("parsing plain private key failed %v", err)
		}

		return signer, nil
	}
}

func parsePemBlock(block *pem.Block) (interface{}, error) {
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing PKCS private key failed %v", err)
		} else {
			return key, nil
		}
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing EC private key failed %v", err)
		} else {
			return key, nil
		}
	case "DSA PRIVATE KEY":
		key, err := ssh.ParseDSAPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing DSA private key failed %v", err)
		} else {
			return key, nil
		}
	default:
		return nil, fmt.Errorf("parsing private key failed, unsupported key type %q", block.Type)
	}
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
