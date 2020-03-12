package config

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func ElasticConnection(cfg Config) elasticsearch.Config {
	addresses := cfg.GetStrings(`elastic.addresses`)

	elasticcfg := elasticsearch.Config{
		Addresses: addresses,
	}

	return elasticcfg
}
