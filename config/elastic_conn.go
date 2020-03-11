package config

import (
	"github.com/elastic/go-elasticsearch"
)

func ElasticConnection(cfg Config) elasticsearch.Config {
	addresses := cfg.GetStrings(`elastic.addresses`)
	username := cfg.GetString(`elastic.username`)
	password := cfg.GetString(`elastic.password`)

	elasticcfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	}

	return elasticcfg
}
