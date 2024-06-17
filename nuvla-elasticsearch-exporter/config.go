package nuvlaedge_otc_exporter

import (
	"errors"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"net/url"
)

type Config struct {
	ElasticsearchConfig *ElasticSearchConfig         `mapstructure:"elasticsearch"`
	QueueConfig         exporterhelper.QueueSettings `mapstructure:"sending_queue"`
	RetryConfig         configretry.BackOffConfig    `mapstructure:"retry_on_failure"`
}

type ElasticSearchConfig struct {
	Enabled             bool     `mapstructure:"enabled"`
	Endpoint            string   `mapstructure:"endpoint"`
	Insecure            bool     `mapstructure:"insecure"`
	CaFile              string   `mapstructure:"ca_file"`
	IndexPrefix         string   `mapstructure:"index_prefix"`
	MetricsTobeExported []string `mapstructure:"metrics"`
}

func (cfg *ElasticSearchConfig) ValidateElasticSearch() error {
	if cfg.Endpoint == "" {
		return errors.New("endpoint must be specified")
	}
	if !cfg.Insecure && cfg.CaFile == "" {
		return errors.New("need to give the ca_file if we want to secure connection")
	}
	return validateEndpoint(cfg.Endpoint)
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
	if cfg.ElasticsearchConfig.Enabled {
		err = cfg.ElasticsearchConfig.ValidateElasticSearch()
		if err != nil {
			return err
		}
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
