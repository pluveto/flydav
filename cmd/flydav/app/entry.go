package app

import (
	"fmt"
	"net/http"

	"github.com/pluveto/flydav/cmd/flydav/conf"
	"github.com/pluveto/flydav/cmd/flydav/service"
)

func Run(conf conf.Conf) {

	if len(conf.Auth.User) == 1 {
		fmt.Println("Username:            ", conf.Auth.User[0].Username)
		fmt.Println("Password(Encrypted): ", conf.Auth.User[0].PasswordHash)
	}
	fmt.Println("Address:             ", fmt.Sprintf("http://%s:%d%s", conf.Server.Host, conf.Server.Port, conf.Server.Path))
	fmt.Println("Filesystem:          ", conf.Server.FsDir)

	server := NewWebdavServer(
		service.NewBasicAuthService(conf.Auth.User),
		conf.Server.Host, conf.Server.Port, conf.Server.Path, conf.Server.FsDir,
	)
	if conf.UI.Enabled {
		http.Handle(conf.UI.Path, http.FileServer(http.Dir(conf.UI.Source)))
	}
	server.Listen()
}
