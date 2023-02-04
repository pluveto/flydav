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

	server := &WebdavServer{
		AuthService: &service.BasicAuthService{Users: conf.Auth.User},
		Host:        conf.Server.Host,
		Port:        conf.Server.Port,
	}
	server.Listen()
}
