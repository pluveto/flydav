package main

import (
    "flydav/cmd/flydav/internal/hub"
    "flydav/cmd/flydav/internal/config"
    "flydav/cmd/flydav/internal/logger"
    "os"
)

func main() {
    // 加载配置文件
    cfg, err := config.Load("config.yaml")
    if err != nil {
        logger.Fatal("Error loading config: ", err)
        os.Exit(1)
    }

    // 初始化日志系统
    logger.Init(cfg)

    // 初始化网络监听
    listener, err := hub.NewListener(cfg.Hub)
    if err != nil {
        logger.Fatal("Error creating listener: ", err)
        os.Exit(1)
    }

    // 启动核心API服务
    go core.Start(cfg.Services.Core)

    // 如果启用了WebDAV服务，启动它
    if cfg.Services.Webdav.Enabled {
        go webdav.Start(cfg.Services.Webdav)
    }

    // 启动HTTP索引服务
    go http_index.Start(cfg.Services.HTTPIndex)

    // 如果启用了UI服务，启动它
    if cfg.Services.UI.Enabled {
        go ui.Start(cfg.Services.UI)
    }

    // 启动认证服务
    go auth.Start(cfg.Services.Auth)

    // 阻塞主线程，等待服务运行
    select {}
}
