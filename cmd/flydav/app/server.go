package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pluveto/flydav/pkg/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/webdav"
)

type AuthService interface {
	Authenticate(username, password string) error
	GetAuthorizedSubDir(username string) (string, error)
	GetPathPrefix(username string) (string, error)
}

type WebdavServer struct {
	AuthService AuthService
	Host        string
	Port        int
	Path        string
	FsDir       string
}

func NewWebdavServer(authService AuthService, host string, port int, path string, fsDir string) *WebdavServer {
	return &WebdavServer{
		AuthService: authService,
		Host:        host,
		Port:        port,
		Path:        path,
		FsDir:       fsDir,
	}
}

func (s *WebdavServer) check() {
	if nil == s.AuthService {
		logger.Fatal("AuthService is nil")
	}
	if s.FsDir == "" {
		logger.Fatal("FsDir is empty")
	}
	// if !path.IsAbs(s.FsDir) {
	// logger.Fatal("FsDir is not an absolute path")
	// }

	var err error
	s.FsDir, err = filepath.Abs(s.FsDir)
	if err != nil {
		logger.Fatal("FsDir is not a valid path", err)
	}

	// must exists
	if _, err := os.Stat(s.FsDir); os.IsNotExist(err) {
		logger.Fatal("FsDir does not exist", err)
	}

	if !strings.HasPrefix(s.FsDir, "/home/") {
		webdavTmpDir := filepath.Join(os.TempDir(), "flydav")
		if !strings.HasPrefix(s.FsDir, webdavTmpDir) {
			logger.Warn("You're using a path which isn't under /home/ as mapped directory. This may cause security issues.")
		}
	}
	logger.Debug("FsDir: ", s.FsDir)
}

func (s *WebdavServer) Listen() {
	s.check()

	lock := webdav.NewMemLS()
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
		subFsDir, err := s.AuthService.GetAuthorizedSubDir(username)
		if err != nil {
			http.Error(w, "Internal Error.", http.StatusInternalServerError)
			logger.Errorf("Error when getting authorized sub dir for user %s: %s", username, err)
			return
		}
		userPrefix, err := s.AuthService.GetPathPrefix(username)
		if err != nil {
			http.Error(w, "Internal Error.", http.StatusInternalServerError)
			logger.Errorf("Error when getting path prefix for user %s: %s", username, err)
		}
		davHandler := &webdav.Handler{
			Prefix:     buildPathPrefix(s.Path, userPrefix),
			FileSystem: buildDirName(s.FsDir, subFsDir),
			LockSystem: lock,
			Logger:     davLogger,
		}

		davHandler.ServeHTTP(w, r)
	})
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	err := http.ListenAndServe(addr, nil)
	logger.Fatal("failed to listen and serve on", addr, ":", err)
}

func buildDirName(fsDir, subFsDir string) webdav.Dir {
	if subFsDir == "" {
		return webdav.Dir(fsDir)
	}
	return webdav.Dir(filepath.Join(fsDir, subFsDir))
}

func buildPathPrefix(path, userPrefix string) string {
	if userPrefix == "" {
		return path
	}
	return filepath.Join(path, userPrefix)
}

func davLogger(r *http.Request, err error) {
	ent := logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	})
	if err != nil {
		ent.Error(err)
	} else {
		ent.Info()
	}
}
