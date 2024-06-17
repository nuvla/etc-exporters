package nuvla_elasticsearch_exporter

import (
	"errors"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"net/url"
)

type Config struct {
	Enabled             bool                         `mapstructure:"enabled"`
	Endpoint            string                       `mapstructure:"endpoint"`
	Insecure            bool                         `mapstructure:"insecure"`
	CaFile              string                       `mapstructure:"ca_file"`
	IndexPrefix         string                       `mapstructure:"index_prefix"`
	MetricsTobeExported []string                     `mapstructure:"metrics"`
	QueueConfig         exporterhelper.QueueSettings `mapstructure:"sending_queue"`
	RetryConfig         configretry.BackOffConfig    `mapstructure:"retry_on_failure"`
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
		if cfg.Endpoint == "" {
			return errors.New("endpoint must be specified")
		}
		if !cfg.Insecure && cfg.CaFile == "" {
			return errors.New("need to give the ca_file if we want to secure connection")
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
