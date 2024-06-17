package nuvla_api_exporter

import (
	"errors"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"net/url"
)

type Config struct {
	Enabled     bool                         `mapstructure:"enabled"`
	Insecure    bool                         `mapstructure:"insecure"`
	Endpoint    string                       `mapstructure:"endpoint"`
	ApiKey      string                       `mapstructure:"api_key"`
	ApiSecret   string                       `mapstructure:"api_secret"`
	ResourceId  string                       `mapstructure:"resource_id"`
	QueueConfig exporterhelper.QueueSettings `mapstructure:"sending_queue"`
	RetryConfig configretry.BackOffConfig    `mapstructure:"retry_on_failure"`
}

func validateEndpoint(endpoint string) error {
	if endpoint == "" {
		return errors.New("endpoint must be specified")
	}
	_, err := url.ParseRequestURI(endpoint)
	return err
}

func (cfg *Config) Validate() error {
	var err error
	if cfg.Enabled {
		if cfg.Endpoint == "" || cfg.ApiKey == "" || cfg.ApiSecret == "" || cfg.ResourceId == "" {
			return errors.New("endpoint, api_key, api_secret and resource id must be specified")
		}

		return validateEndpoint(cfg.Endpoint)
	}

	err = cfg.QueueConfig.Validate()
	if err != nil {
		return err
	}
	err = cfg.RetryConfig.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) Unmarshal(conf *confmap.Conf) error {
	return conf.Unmarshal(cfg)
}

var _ component.Config = (*Config)(nil)
var _ confmap.Unmarshaler = (*Config)(nil)
