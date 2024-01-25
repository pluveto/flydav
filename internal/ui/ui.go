// internal/ui/ui.go
package ui

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pluveto/flydav/internal/config"
	"github.com/pluveto/flydav/internal/logger"
)

type UIModule struct {
	Config config.UIConfig
	// 可能还有其他的状态或依赖
}

func NewUIModule(cfg config.UIConfig) *UIModule {
	return &UIModule{
		Config: cfg,
	}
}

func (uis *UIModule) RegisterRoutes(router *mux.Router) {
	router.PathPrefix("/app").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

func (uis *UIModule) Start() error {
	logger.Info("Starting UI Module")
	return nil
}
