package remoteexec

import (
	"fmt"

	"github.com/lrwx00t/go-remoteexec/config"
	"github.com/lrwx00t/go-remoteexec/ssh_utils"
)

func ServerClone(ge config.GlobalExec) {
	// clone_repo_cmd := "git clone https://github.com/0xack13/go-bootstrap"
	clone_repo_cmd := fmt.Sprintf("git clone %s", ge.Config.Repo) + "\n"
	ssh_utils.Session_execute(clone_repo_cmd, ge.Connection)
}
