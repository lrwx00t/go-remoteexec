package remoteexec

import (
	"fmt"
	"net/url"
	"path"

	"github.com/lrwx00t/go-remoteexec/config"
	"github.com/lrwx00t/go-remoteexec/ssh_utils"
)

func Cleanup_exec(ge config.GlobalExec) {
	u, err := url.Parse(ge.Config.Repo)
	if err != nil {
		panic(err)
	}

	// Extract the name of the repository
	repoName := path.Base(u.Path)
	delete_cmd := fmt.Sprintf("rm -fR %s", repoName)
	ssh_utils.Session_execute(delete_cmd, ge.Connection)
}
