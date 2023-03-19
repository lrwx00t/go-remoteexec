# go-remoteexec
Bootstrap servers using Go Remote Execution via SSH

Usage
To run the bootstrap scripts on an instance with password:
```bash
REMOTEEXEC_PASS="MYPASS" go run main.go -a=159.89.50.2
```
Using private keys is also supported with -i identity option:
```bash
go run main.go -a=159.89.50.2 -i=~/.ssh/id_rsa_droplets
```
