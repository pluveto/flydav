// internal/hub/listener.go
package hub

import (
    "net"
    "flydav/cmd/flydav/internal/config"
)

func NewListener(cfg config.HubConfig) (net.Listener, error) {
    // 根据配置决定是创建普通的TCP监听还是TLS监听
    // ...
    return listener, nil
}
