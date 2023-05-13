package remoteexec

import (
	"fmt"

	"github.com/lrwx00t/go-remoteexec/config"
	"github.com/lrwx00t/go-remoteexec/ssh_utils"
)

func Ps1(ps string) string {
	// TODO check if not nil
	return fmt.Sprintf(`echo 'PS1="[\u\H %s ]# "' >> ~/.bashrc`+"\n", ps)
}

func PS1_exec(ge config.GlobalExec) {
	ssh_utils.Session_execute(Ps1(ge.Config.HostAlias), ge.Connection)
}
