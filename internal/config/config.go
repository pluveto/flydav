// internal/config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/pluveto/flydav/internal/hub"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Hub      hub.HubConfig  `yaml:"hub"`
	Services ServicesConfig `yaml:"services"`
	Log      LogConfig      `yaml:"log"`
}

type ServicesConfig struct {
	Core      CoreConfig      `yaml:"core"`
	WebDAV    WebDAVConfig    `yaml:"webdav"`
	HTTPIndex HTTPIndexConfig `yaml:"http_index"`
	UI        UIConfig        `yaml:"ui"`
	Auth      AuthConfig      `yaml:"auth"`
}

type CORSConfig struct {
	Enabled bool     `yaml:"enabled"`
	Origins []string `yaml:"origins"`
}

type LogConfig struct {
	Enabled bool   `yaml:"enabled"`
	Level   string `yaml:"level"`
	Path    string `yaml:"path"`
}

type WebDAVConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type HTTPIndexConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type UIConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %s", err)
	}

	return config, nil
}
