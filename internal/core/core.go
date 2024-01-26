// internal/core/core.go
package core

import (
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
	// 实现注册路由的逻辑
}

func (cm *CoreModule) Start() error {
	logger.Info("starting core module")
	return nil
}
