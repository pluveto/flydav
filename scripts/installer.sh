#!/usr/bin/env bash
#
# This script will download latest release from
#   https://github.com/pluveto/flydav/releases
# and install it to /usr/local/bin/flydav
# with configuration file /etc/flydav/flydav.toml

DOWNLOADER="curl"

# ------------------ common ------------------

get_downloader() {
    if [ -x "$(command -v curl)" ]; then
        DOWNLOADER="curl"
    elif [ -x "$(command -v wget)" ]; then
        DOWNLOADER="wget"
    else
        echo "No downloader found. Please install curl or wget."
        exit 1
    fi
}

must_has() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "This script requires $1 but it's not installed. Aborting."
        exit 1
    fi
}

must_run() {
    if ! "$@"; then
        echo "Failed to run $*"
        exit 1
    fi
}

must_be_root() {
    if [ "$(id -u)" != "0" ]; then
        echo "This script must be run as root"
        exit 1
    fi
}

must_be_linux() {
    if [ "$(uname)" != "Linux" ]; then
        echo "This script only works on Linux"
        exit 1
    fi
}


supported_platforms=(
    linux-386
    linux-amd64
    linux-arm
    linux-arm64
    mac-amd64
    mac-arm64
)

# ------------------ adhoc ------------------

REPO_API=https://api.github.com/repos/pluveto/flydav
REPO=https://github.com/pluveto/flydav
VERSION=
DEBUG=1

get_latest_release() {
    echo "Getting latest release from $REPO_API"
    if [ "$DOWNLOADER" = "curl" ]; then
        VERSION=$(curl -s "$REPO_API/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif [ "$DOWNLOADER" = "wget" ]; then
        VERSION=$(wget -qO- "$REPO_API/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    fi
    if [ -z "$VERSION" ]; then
        echo "Failed to get latest release"
        exit 1
    fi
}

clean_up() {
    if [ -d /tmp/flydav_install ]; then
        rm -rfd /tmp/flydav_install
    fi
}


must_has unzip


get_downloader

get_latest_release

must_be_root

must_be_linux

ARCH=$(uname -m)
ARCH=${ARCH/x86_64/amd64}
ARCH=${ARCH/aarch64/arm64}

DOWNLOAD_URL="$REPO/releases/download/$VERSION/flydav-app-linux-$ARCH.zip"

echo "Downloading flydav $VERSION from $DOWNLOAD_URL"

must_run mkdir -p /tmp/flydav_install

if [ ! -f /tmp/flydav_install/flydav.zip ]; then
    if [ "$DOWNLOADER" = "curl" ]; then
        must_run curl -L "$DOWNLOAD_URL" -o /tmp/flydav_install/flydav.zip
    elif [ "$DOWNLOADER" = "wget" ]; then
        must_run wget -O /tmp/flydav_install/flydav.zip "$DOWNLOAD_URL"
    fi
    echo "Downloaded flydav to /tmp/flydav_install/flydav.zip"
else
    echo "Skip downloading flydav"
fi

# assert flydav.zip size is larger than 1MB
if [ "$(stat -c%s /tmp/flydav_install/flydav.zip)" -lt 1000000 ]; then
    echo "Downloaded flydav is too small, please check your network"
    clean_up
    exit 1
fi

must_run unzip /tmp/flydav_install/flydav.zip -d /tmp/flydav_install

must_run chmod +x /tmp/flydav_install/flydav


DOWNLOAD_CONFIG=1
# if already has configuration, ask user to keep it or not
if [ -f /etc/flydav/flydav.toml ]; then
    echo "Configuration file /etc/flydav/flydav.toml already exists."
    echo "Do you want to keep it? (y/n)"
    read -r answer
    if [ "$answer" = "n" ]; then
        must_run rm -f /etc/flydav/flydav.toml
    else
        DOWNLOAD_CONFIG=0
    fi
fi

if [ "$DOWNLOAD_CONFIG" = "1" ]; then

    TMP_CONFIG_PATH=/tmp/flydav_install/flydav.toml
    CONFIG_DOWNLOAD_URL="https://raw.githubusercontent.com/pluveto/flydav/main/conf/config.default.toml"

    echo "Downloading configuration file from $CONFIG_DOWNLOAD_URL"

    if [ "$DOWNLOADER" = "curl" ]; then
        must_run curl -L "$CONFIG_DOWNLOAD_URL" -o "$TMP_CONFIG_PATH"
    elif [ "$DOWNLOADER" = "wget" ]; then
        must_run wget -O "$TMP_CONFIG_PATH" "$CONFIG_DOWNLOAD_URL"
    fi
    echo "Downloaded configuration file to $TMP_CONFIG_PATH"

    read -r -d '' USER_CONFIG_TMPL <<'EOF'
    [[auth.user]]
    username = "::::username::::"
    sub_fs_dir = "::::sub_fs_dir::::"
    sub_path = "::::sub_path::::"
    password_hash = "::::password_hash::::"
    password_crypt = "sha256"
EOF

    USERNAME=flydav

    echo "Username (default: flydav):"
    read -r answer
    if [ -n "$answer" ]; then
        USERNAME="$answer"
    fi

    echo "Password for $USERNAME:"
    read -r PASSWORD
    # password must be longer than or eq 10 characters
    while [ ${#PASSWORD} -lt 10 ]; do
        echo "Password must be longer than or eq 10 characters"
        echo "Please enter a password for $USERNAME:"
        read -r PASSWORD
    done
    PASSWORD_HASH=$(echo -n "$PASSWORD" | sha256sum | cut -d ' ' -f 1)

    FS_DIR=/tmp/flydav
    SUB_FS_DIR=
    echo "Filesystem root dir: (default: /tmp/flydav)"
    read -r answer
    if [ -n "$answer" ]; then
        FS_DIR="$answer"
    fi
    echo "Sub filesystem dir: (default: $SUB_FS_DIR)"
    read -r answer
    if [ -n "$answer" ]; then
        SUB_FS_DIR="$answer"
    fi
    TMP_USER_FS_ROOT="$FS_DIR/$SUB_FS_DIR"
    echo "User dir will be $TMP_USER_FS_ROOT Do you want to continue? (y/n)"
    read -r answer
    
    if [ "$answer" != "y" ] && [ "$answer" != "yes" ] && [ -n "$answer" ]; then
        echo "Installation cancelled"
        clean_up
        exit 1
    fi

    HTTP_HOST=0.0.0.0
    HTTP_PORT=7086
    echo "HTTP host(default: $HTTP_HOST): "
    read -r answer
    if [ -n "$answer" ]; then
        HTTP_HOST="$answer"
    fi
    echo "HTTP port(default: $HTTP_PORT): "
    read -r answer
    if [ -n "$answer" ]; then
        HTTP_PORT="$answer"
    fi

    HTTP_PATH_PREFIX=/webdav
    echo "HTTP path prefix(default: $HTTP_PATH_PREFIX): "
    read -r answer
    if [ -n "$answer" ]; then
        HTTP_PATH_PREFIX="$answer"
    fi

    echo "The service will be running at http://$HTTP_HOST:$HTTP_PORT$HTTP_PATH_PREFIX"

    echo "Do you want to continue? (y/n)"
    read -r answer
    if [ "$answer" != "y" ] && [ "$answer" != "yes" ] && [ -n "$answer" ]; then
        echo "Installation cancelled"
        clean_up
        exit 1
    fi

    echo "Creating configuration file"
    must_run mkdir -p /etc/flydav
 
    must_run sed -i "/# add more users here/r /dev/stdin" "$TMP_CONFIG_PATH" <<< "$USER_CONFIG_TMPL"

    function sedeasy {
      sed -i "s/$(printf '%s\n' "$1" | sed -e 's/\([[\/.*]\|\]\)/\\&/g')/$(echo $2 | sed -e 's/[\/&]/\\&/g')/g" $3
    }
    must_run sedeasy "::::username::::" "$USERNAME" "$TMP_CONFIG_PATH"
    must_run sedeasy "::::sub_fs_dir::::" "$SUB_FS_DIR" "$TMP_CONFIG_PATH"
    must_run sedeasy "::::sub_path::::" "$SUB_PATH" "$TMP_CONFIG_PATH"
    must_run sedeasy "::::password_hash::::" "$PASSWORD_HASH" "$TMP_CONFIG_PATH"

    must_run sedeasy "host = \"0.0.0.0\"" "host = \"$HTTP_HOST\"" "$TMP_CONFIG_PATH"
    must_run sedeasy "port = 7086" "port = $HTTP_PORT" "$TMP_CONFIG_PATH"
    must_run sedeasy "path = \"\/webdav\"" "path = \"$HTTP_PATH_PREFIX\"" "$TMP_CONFIG_PATH"
    must_run sedeasy "fs_dir = \"\/tmp\/flydav\"" "fs_dir = \"$FS_DIR\"" "$TMP_CONFIG_PATH"

    echo "Configuration file created at $TMP_CONFIG_PATH"

    must_run mv "$TMP_CONFIG_PATH" /etc/flydav/flydav.toml
fi

echo "Installing binary"
must_run mv "/tmp/flydav_install/flydav" "/usr/local/bin/flydav"

# check if service is already installed
if [ -f /etc/systemd/system/flydav.service ]; then
    echo "Service already installed"
    echo "Do you want to reinstall? (y/n)"
    read -r answer
    if [ "$answer" != "y" ] && [ "$answer" != "yes" ] && [ -n "$answer" ]; then
        echo "Installation cancelled"
        clean_up
        exit 1
    fi

    echo "Removing old service"
    must_run systemctl stop flydav
    must_run systemctl disable flydav
    must_run rm -r /etc/systemd/system/flydav.service
    must_run systemctl daemon-reload

fi

echo "Creating user flydav"
# create if not exists
if ! id -u flydav >/dev/null 2>&1; then
    must_run useradd -r -s /bin/false flydav
fi

# if TMP_USER_FS_ROOT is not existing
if [ ! -d "$TMP_USER_FS_ROOT" ]; then

    echo "Creating directory $TMP_USER_FS_ROOT"
    must_run mkdir -p "$TMP_USER_FS_ROOT"

    echo "Setting permissions for $TMP_USER_FS_ROOT"
    must_run chown -R flydav:flydav "$TMP_USER_FS_ROOT"

fi

echo "Creating systemd service"

read -r -d '' SERVICE_TMPL <<'EOF'
[Unit]
Description=Flydav WebDAV server
After=network.target

[Service]
User=flydav
Group=flydav
ExecStart=/usr/local/bin/flydav -c /etc/flydav/flydav.toml
Restart=on-failure
RestartSec=5
StartLimitInterval=60s
StartLimitBurst=3

[Install]
WantedBy=multi-user.target
EOF

must_run echo "$SERVICE_TMPL" > /etc/systemd/system/flydav.service

echo "Enabling systemd service"
must_run systemctl daemon-reload
must_run systemctl enable flydav.service

echo "Starting systemd service"
must_run systemctl start flydav.service

echo "Installation complete"
