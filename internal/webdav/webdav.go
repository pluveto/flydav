// internal/webdav/webdav.go
package webdav

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pluveto/flydav/internal/auth"
	"github.com/pluveto/flydav/internal/config"
	"github.com/pluveto/flydav/internal/logger"
	"github.com/pluveto/flydav/pkg/storage"
)

type WebDAVModule struct {
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
	Auth    *auth.AuthModule
}

func NewWebDAVModule(cfg config.WebDAVConfig, storage storage.Storage, auth *auth.AuthModule) *WebDAVModule {
	return &WebDAVModule{
		Config:  cfg,
		Storage: storage,
		Auth:    auth,
	}
}

func (wds *WebDAVModule) RegisterRoutes(router *mux.Router) {
	router.PathPrefix(wds.Config.Path).
		Methods("GET", "POST", "PUT", "DELETE", "MKCOL", "COPY", "MOVE", "OPTIONS").
		HandlerFunc(wds.handleWebDAV)

}

func (wds *WebDAVModule) handleWebDAV(w http.ResponseWriter, r *http.Request) {
	// 检查用户是否已经通过验证
	if !wds.Auth.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	filePath := r.URL.Path

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
		// 你需要解析目标头部（Destination header），然后调用 Storage 的 Move 或者 Copy 方法
		// ...

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

func (wds *WebDAVModule) Start() error {
	logger.Info("starting web_dav module")
	return nil
}
