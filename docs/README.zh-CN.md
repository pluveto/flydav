# Flydav

```text
/flydav
|-- /cmd
|   |-- /flydav
|       |-- main.go         # 程序入口点，解析配置，初始化和启动服务
|-- /internal
|   |-- /hub                # 处理网络监听和TLS配置
|   |-- /core               # 核心API逻辑
|   |-- /webdav             # WebDAV服务实现
|   |-- /http_index         # HTTP索引服务实现
|   |-- /ui                 # 用户界面服务实现（如果启用的话）
|   |-- /auth               # 认证服务实现
|   |-- /config             # 配置加载和解析
|   |-- /logger             # 日志配置和管理
|-- /pkg
|   |-- /api                # 公共API定义和工具
|   |-- /storage            # 存储抽象层，可能包括本地文件系统和S3兼容存储
|   |-- /authenticator      # 认证器接口和实现
|-- go.mod                  # Go模块定义
|-- go.sum                  # Go模块依赖校验和
|-- Dockerfile              # 用于构建Docker镜像的Dockerfile
|-- config.yaml             # 示例配置文件
```
