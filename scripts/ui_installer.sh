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

DOWNLOAD_URL="$REPO/releases/download/$VERSION/flydav-ui-dist.zip"

echo "Downloading flydav $VERSION from $DOWNLOAD_URL"

must_run mkdir -p /tmp/flydav_install

if [ ! -f /tmp/flydav_install/flydav-ui-dist.zip ]; then
    if [ "$DOWNLOADER" = "curl" ]; then
        must_run curl -L "$DOWNLOAD_URL" -o /tmp/flydav_install/flydav-ui-dist.zip
    elif [ "$DOWNLOADER" = "wget" ]; then
        must_run wget -O /tmp/flydav_install/flydav-ui-dist.zip "$DOWNLOAD_URL"
    fi
    echo "Downloaded flydav to /tmp/flydav_install/flydav-ui-dist.zip"
else
    echo "Skip downloading flydav"
fi

# assert flydav-ui-dist.zip size is larger than 10KB
if [ "$(stat -c%s /tmp/flydav_install/flydav-ui-dist.zip)" -lt 10000 ]; then
    echo "Downloaded flydav-ui is too small, please check your network"
    clean_up
    exit 1
fi

echo "Install dir for flydav-ui (default: /usr/local/share/flydav/ui):"
read -r answer
if [ -n "$answer" ]; then
    UI_INSTALL_DIR="$answer"
else
    UI_INSTALL_DIR=/usr/local/share/flydav/ui
fi

mkdir -p "$UI_INSTALL_DIR"

must_run unzip /tmp/flydav_install/flydav-ui-dist.zip -d $UI_INSTALL_DIR

echo "Flydav UI is installed to $UI_INSTALL_DIR"
echo "Edit your flydav config file and restart to enable UI"

clean_up
