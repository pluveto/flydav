package config

type CoreConfig struct {
	Path    string               `yaml:"path"`
	CORS    CORSConfig           `yaml:"cors"`
	Backend StorageBackendConfig `yaml:"backend"`
	Log     LogConfig            `yaml:"log"`
}

type StorageBackendConfig struct {
	Local LocalConfig `yaml:"local"`
	S3    S3Config    `yaml:"s3"`
}

func (sbc *StorageBackendConfig) GetEnabledBackend() string {
	if sbc.Local.Enabled {
		return "local"
	}
	if sbc.S3.Enabled {
		return "s3"
	}
	return ""
}

type LocalConfig struct {
	Enabled bool   `yaml:"enabled"`
	BaseDir string `yaml:"base_dir"`
}

type S3Config struct {
	Enabled   bool   `yaml:"enabled"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
}
