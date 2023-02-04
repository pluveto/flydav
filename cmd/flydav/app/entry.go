package app

import (
	"fmt"

	"github.com/pluveto/flydav/cmd/flydav/conf"
	"github.com/pluveto/flydav/cmd/flydav/service"
)

func Run(conf conf.Conf) {

	fmt.Println("Serving on:          ", fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))
	fmt.Println("Username:            ", conf.Auth.User[0].Username)
	// fmt.Println("Password(Encrypted): ", conf.Auth.User[0].PasswordHash)

	server := NewWebdavServer(
		service.NewBasicAuthService(conf.Auth.User),
		conf.Server.Host, conf.Server.Port, conf.Server.Path, conf.Server.FsDir,
	)
	server.Listen()
}
