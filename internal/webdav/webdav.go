// internal/webdav/webdav.go
package webdav

import (
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pluveto/flydav/internal/config"
	"github.com/pluveto/flydav/internal/logger"
	"github.com/pluveto/flydav/pkg/authenticator"
	"github.com/pluveto/flydav/pkg/storage"
)

type WebDAVService struct {
	Config config.WebDAVConfig
	/*
		type Storage interface {
			WriteAll(path string, data []byte) error
			ReadAll(path string) ([]byte, error)
			Delete(path string) error
			Write(path string, data []byte, offset int64) error
			Read(path string, offset int64, length int64) ([]byte, error)
			Stat(path string) (Metadata, error)
			List(path string) ([]Metadata, error)
			Size(path string) (int64, error)
			Move(src string, dst string) error
			Merge(srcs []string, dst string) error
		}
	*/
	Storage storage.Storage
	Auth    authenticator.Authenticator
}

func NewWebDAVService(cfg config.WebDAVConfig, storage storage.Storage, auth authenticator.Authenticator) *WebDAVService {
	return &WebDAVService{
		Config:  cfg,
		Storage: storage,
		Auth:    auth,
	}
}

func (wds *WebDAVService) RegisterRoutes(router *mux.Router) {
	router.PathPrefix(wds.Config.Path).
		Methods("GET", "POST", "PUT", "DELETE", "MKCOL", "COPY", "MOVE", "OPTIONS").
		HandlerFunc(wds.handleWebDAV)

}

func (wds *WebDAVService) handleWebDAV(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="WebDAV"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authenticated, authErr := wds.Auth.Authenticate(username, password)
	if !authenticated || authErr != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		if authErr != nil {
			logger.Error("error authenticating user: ", authErr)
		}
		return
	}

	// 获取用户的根目录
	rootDir := wds.Auth.GetRootDir(username)
	// 确保文件路径是以根目录开始的
	filePath := rootDir + r.URL.Path
	// 检查用户是否有权访问该路径
	hasPermission, err := wds.Auth.Authorize(username, filePath, config.PermissionRead) // or PermissionWrite, depending on the method
	if err != nil || !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	switch r.Method {
	case "GET":
		data, err := wds.Storage.ReadAll(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)

	case "POST":
		// POST 通常不用于传统的 WebDAV 操作，可能需要根据你的具体需求来实现

	case "PUT":
		data, err := readRequestBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = wds.Storage.WriteAll(filePath, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "DELETE":
		err := wds.Storage.Delete(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "MKCOL":
		// 创建一个新的集合（目录）
		err := wds.Storage.CreateDirectory(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "COPY", "MOVE":
		destination := r.Header.Get("Destination")
		if destination == "" {
			http.Error(w, "Destination header missing", http.StatusBadRequest)
			return
		}

		// 解析目标路径，并确保它以用户的根目录开始
		destPath := rootDir + destination
		// 检查用户是否有权访问目标路径
		hasPermission, err := wds.Auth.Authorize(username, destPath, config.PermissionWrite)
		if err != nil || !hasPermission {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if r.Method == "COPY" {
			_, err = wds.Storage.Copy(filePath, destPath)
		} else { // MOVE
			err = wds.Storage.Move(filePath, destPath)
		}

		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, "Not Found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

	case "OPTIONS":
		// 返回支持的方法
		w.Header().Set("Allow", "OPTIONS, GET, POST, PUT, DELETE, MKCOL, COPY, MOVE")
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// 读取请求体的辅助函数
func readRequestBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

func (wds *WebDAVService) Start() error {
	logger.Info("starting web_dav module")
	return nil
}
