package app

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/pluveto/flydav/pkg/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/webdav"
)

type AuthService interface {
	Authenticate(username, password string) error
}

type WebdavServer struct {
	AuthService AuthService
	Host        string
	Port        int
	MappedDir   string
}

func (s *WebdavServer) Listen() {
	if nil == s.AuthService {
		logger.Fatal("AuthService is nil")
	}
	if s.MappedDir == "" {
		logger.Fatal("MappedDir is empty")
	}
	if !path.IsAbs(s.MappedDir) {
		logger.Fatal("MappedDir is not an absolute path")
	}
	if !strings.HasPrefix(s.MappedDir, "/home/") {
		logger.Warn("You're using a path which isn't under /home/ as mapped directory. This may cause security issues.")
	}

	davfs := &webdav.Handler{
		FileSystem: webdav.Dir(s.MappedDir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			ent := logger.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
			})
			if err != nil {
				ent.Error(err)
			} else {
				ent.Info()
			}
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		err := s.AuthService.Authenticate(username, password)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			logger.Error("Unauthorized: ", err)
			return
		}
		davfs.ServeHTTP(w, r)
	})
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	err := http.ListenAndServe(addr, nil)
	logger.Fatal("failed to listen and serve on", addr, ":", err)
}
