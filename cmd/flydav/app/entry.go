package app

import (
	"fmt"
	"net/http"
	"strings"

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
	if conf.CORS.Enabled {
		server.AddMiddleware(func(next http.HandlerFunc) http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", strings.Join(conf.CORS.AllowedOrigins, ","))
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(conf.CORS.AllowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(conf.CORS.AllowedHeaders, ","))
				w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprintf("%v", conf.CORS.AllowCredentials))
				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}
				if conf.CORS.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", conf.CORS.MaxAge))
				}
				if len(conf.CORS.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(conf.CORS.ExposedHeaders, ","))
				}

				next.ServeHTTP(w, r)
			})
		})
	}
	server.Listen()
}
