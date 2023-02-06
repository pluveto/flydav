# Flydav

FlyDav 是一个轻量级的开源 webdav 服务器，它提供的一些核心功能可以满足个人用户和组织的需求。

FlyDav 的文件大小很小，对于需要快速高效的 webdav 服务器的用户来说，它是理想的解决方案。它提供基本认证，支持多用户，并允许每个用户拥有不同的根目录和路径前缀，这些都是非常好的隔离。此外，FlyDav 还提供了日志轮换和密码安全，是需要安全可靠的 webdav 解决方案的用户的最佳选择。

FlyDav 的目标是保持方便和简洁，用户能在短时间内将服务部署起来。因此任何无关紧要的功能都不会被加入，从而避免臃肿。

## 在 30 秒内开始使用

1. 首先从 [发布页](https://github.com/pluveto/flydav/releases) 下载 FlyDav。
2. 运行 `./flydav -H 0.0.0.0` 来启动服务器。然后你要输入默认用户 `flydav` 的密码。
3. 在你的 webdav 客户端（比如 RaiDrive）中打开 `http://YOUR_IP:7086/webdav`。

## 命令行选项

```bash
$ flydav -h
--------------------------------------------------------------------------------
用法: flydav [--host HOST] [--port PORT] [--user USER] [--verbose] [--config CONFIG] 。

选项:
  --host HOST, -H HOST 主机地址
  --port PORT, -p PORT 端口
  --user USER, -u USER 用户名
  --verbose, -v   启用详细输出输出（如果你打算报告或调试错误，这将非常有用）
  --config CONFIG, -c CONFIG
                         配置文件
  --help, -h 显示此帮助并退出
```

如果你有一个配置文件，你可以忽略这些命令行选项。运行 `flydav -c /path/to/config.toml` 来启动服务器。

如果你想用主机、端口、用户名和一次性密码快速启动服务器，你可以运行 `flydav -H IP -p PORT -u USERNAME` 来启动服务器。然后你再输入用户的密码。然后服务器将在 `http://IP:PORT/` 提供服务。

## 配置 FlyDav

尽管 FlyDav 提供了一些命令行选项，但是你也可以使用配置文件来配置它。从而避免在每次启动服务器时都输入密码。

1. 下载 FlyDav。
2. 现在你有了这个软件，你需要为它创建一个配置文件。首先创建一个名为 "flydav.toml" 的新文件。
3. 在配置文件中，你将需要添加以下信息。
    - `[服务器]`。这一部分将定义 webdav 服务器的主机、端口和路径。
    - `host`: 主机的 IP 地址。如果你想让任何 IP 地址都能访问该服务器，这应该被设置为 "0.0.0.0"。
    - `port`: webdav 服务器要使用的端口号。
    - `path`: webdav 服务器的路径。
    - `fs_dir`: 服务器上存放 webdav 文件的目录。
    - `[auth]`: 这一部分将定义 webdav 服务器的认证设置。
    - `[[auth.user]]`: 这一节将为每个可以访问 webdav 服务器的用户定义用户名和凭证。
        - `username`: 用户的用户名。
        - `sub_fs_dir': 用户可以访问的 fs_dir 的子目录。
        - `sub_path`: 用户访问 webdav 服务器的路径
        - `password_hash`: 用户的散列密码。
        - `password_crypt`: 用于哈希密码的哈希算法的类型。这应该被设置为 "bcrypt"。
    - `[log]`: 这一部分将定义 webdav 服务器的日志设置。
    - `level`: 服务器的日志级别。这可以设置为 "debug"、"info"、"warning"、"error" 或 "fatal"。
    - `[[log.file]]`。这个小节将定义日志文件的设置。如果你不想将日志记录到一个文件中，请忽略这个小节。
        - `format`: 日志文件的格式。这可以设置为 "json" 或 "text"。
        - `path`: 日志文件的路径。
        - `max_size`: 日志文件的最大尺寸，以兆字节为单位。
        - `max_age`: 日志文件的最大年龄，以天为单位。
    - `[[log.stdout]]`: 本小节将定义日志输出到控制台的设置。如果你不想向控制台输出日志，请忽略这个小节。
        - `format`: 日志输出的格式。可以设置为 "json" 或 "text"。
        - `output`: 日志输出的输出流。可以设置为 "stdout" 或 "stderr"。

4. 保存配置文件并运行 FlyDav 服务器。现在你应该可以用配置好的设置访问 webdav 服务器了。

要获得一个配置文件的例子，请到 [conf dir](https://github.com/pluveto/flydav/blob/main/conf)。

## 以服务方式安装

### 在 Linux 上以服务方式安装

1. 创建编辑 `/etc/systemd/system/flydav.service`，并添加以下内容：

```ini
[Unit]
Description = Flydav Server
After = network.target syslog.target
Wants = network.target

[Service]
Type = simple
# !!! 把配置文件和程序位置改成你自己的 !!!
ExecStart = /usr/bin/flydav -c /etc/flydav/flydav.toml

[Install]
WantedBy = multi-user.target
```

2. 运行 `systemctl daemon-reload`，重新加载systemd守护程序。
3. 运行 `systemctl enable flydav` 来启用该服务。
4. 运行 `systemctl start flydav` 来启动服务。

### 管理该服务

- 运行 `systemctl status flydav` 来检查服务的状态。
- 运行 `systemctl stop flydav` 来停止该服务。

## 功能

- [x] 基本认证
- [x] 多个用户
- [x] 每个用户的根目录不同
- [x] 每个用户有不同的路径前缀
- [x] 日志
- [ ] SSL 
    - 正在支持
    - 你可以使用 Nginx 这样的反向代理来启用 SSL

## 许可证

在MIT许可下授权——详见[LICENSE](../LICENSE)文件
