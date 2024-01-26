// internal/core/core.go
package core

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pluveto/flydav/internal/config"
	"github.com/pluveto/flydav/internal/logger"
)

type CoreModule struct {
	Config config.CoreConfig
	// 可能还有其他的状态或依赖
}

func NewCoreModule(cfg config.CoreConfig) *CoreModule {
	return &CoreModule{
		Config: cfg,
	}
}

func (cm *CoreModule) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/internal/health", cm.handleHealth).Methods("GET")
	router.HandleFunc("/", cm.handleRoot).Methods("GET")
}

func (cm *CoreModule) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (cm *CoreModule) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("flydav v1.0"))
}

func (cm *CoreModule) Start() error {
	logger.Info("starting core module")
	return nil
}
