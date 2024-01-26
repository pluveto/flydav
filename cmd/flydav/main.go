package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pluveto/flydav/internal/auth"
	"github.com/pluveto/flydav/internal/config"
	"github.com/pluveto/flydav/internal/core"
	"github.com/pluveto/flydav/internal/http_index"
	"github.com/pluveto/flydav/internal/hub"
	"github.com/pluveto/flydav/internal/logger"
	"github.com/pluveto/flydav/internal/ui"
	"github.com/pluveto/flydav/internal/webdav"
	"github.com/pluveto/flydav/pkg/storage"
)

func main() {
	// 加载配置文件
	cfg, err := config.Load("config.yaml")
	if err != nil {
		logger.Fatal("error loading config: ", err)
		os.Exit(1)
	}

	// 初始化日志系统
	logger.Init(cfg.Log)
	logger.Info("starting flydav")

	// 初始化网络监听
	listener, err := hub.NewListener(cfg.Hub)
	if err != nil {
		logger.Fatal("error creating listener: ", err)
		os.Exit(1)
	}
	router := mux.NewRouter()
	storage := storage.NewStorage(cfg.Services.Core.Backend)
	authModule := auth.NewAuthModule(cfg.Services.Auth)

	coreModule := core.NewCoreModule(cfg.Services.Core)
	coreModule.RegisterRoutes(router)

	webdavModule := webdav.NewWebDAVService(cfg.Services.WebDAV, storage, authModule)
	webdavModule.RegisterRoutes(router)

	httpIndexModule := http_index.NewHTTPIndexModule(cfg.Services.HTTPIndex, storage, authModule)
	httpIndexModule.RegisterRoutes(router)

	uiModule := ui.NewUIModule(cfg.Services.UI)
	uiModule.RegisterRoutes(router)

	authModule.RegisterRoutes(router)

	go authModule.Start()
	go uiModule.Start()
	go httpIndexModule.Start()
	go webdavModule.Start()
	go coreModule.Start()

	printRoutes(router)

	logger.Info("listening on ", cfg.Hub.GetListenAddress())
	log.Fatal(http.Serve(listener, router))

}

func printRoutes(router *mux.Router) {
	s := ""
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			logger.Error(err)
		}
		methods, err := route.GetMethods()
		if err != nil {
			logger.Error(err)
		}

		method := ""
		for _, m := range methods {
			method += m + " "
		}

		s += " - " + method + path + "\n"
		return nil
	})

	logger.Info("\nregistered_routes: \n", s)
}
