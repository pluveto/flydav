PROJECT_NAME=flydav-app
BIN_NAMES=flydav
GOARCHS=amd64 386 arm arm64
GOARCHS_MAC=amd64 arm64

dev: linux

default: all

all: windows linux mac

prepare:
	@mkdir -p dist

windows: prepare
	for BIN_NAME in $(BIN_NAMES); do \
		[ -z "$$BIN_NAME" ] && continue; \
		for GOARCH in $(GOARCHS); do \
			mkdir -p dist/windows_$$GOARCH; \
			GOOS=windows GOARCH=$$GOARCH go build -o dist/windows_$$GOARCH/$$BIN_NAME.exe cmd/$$BIN_NAME/main.go; \
		done \
	done

linux: prepare
	for BIN_NAME in $(BIN_NAMES); do \
		[ -z "$$BIN_NAME" ] && continue; \
		for GOARCH in $(GOARCHS); do \
			mkdir -p dist/linux_$$GOARCH; \
			GOOS=linux GOARCH=$$GOARCH go build -o dist/linux_$$GOARCH/$$BIN_NAME cmd/$$BIN_NAME/main.go; \
		done \
	done

mac: prepare
	for BIN_NAME in $(BIN_NAMES); do \
		[ -z "$$BIN_NAME" ] && continue; \
		for GOARCH in $(GOARCHS_MAC); do \
			mkdir -p dist/mac_$$GOARCH; \
			GOOS=darwin GOARCH=$$GOARCH go build -o dist/mac_$$GOARCH/$$BIN_NAME cmd/$$BIN_NAME/main.go; \
		done \
	done

package: all
	for GOARCH in $(GOARCHS); do \
		(cd dist/windows_$$GOARCH/ && zip -q -r ../../dist/$(PROJECT_NAME)-windows-$$GOARCH.zip .); \
		(cd dist/linux_$$GOARCH/ && zip -q -r ../../dist/$(PROJECT_NAME)-linux-$$GOARCH.zip .); \
	done

	for GOARCH in $(GOARCHS_MAC); do \
		(cd dist/mac_$$GOARCH/ && zip -q -r ../../dist/$(PROJECT_NAME)-mac-$$GOARCH.zip .); \
	done

	ARCH_RELEASE_DIRS=$$(find dist -type d -name "*_*"); \
	for ARCH_RELEASE_DIR in $$ARCH_RELEASE_DIRS; do \
		cp conf/config.default.toml $$ARCH_RELEASE_DIR/config.toml; \
		rm -rfd $$ARCH_RELEASE_DIR; \
	done

test:
	go test -v ./...

run:
	go run cmd/$(APP)/main.go $(ARGS)

clean:
	rm -rfd dist

.PHONY: all, default, clean
