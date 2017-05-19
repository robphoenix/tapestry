package tapestry

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
)

// Config ...
type Config struct {
	URL           string
	User          string
	Password      string
	DataSrc       string
	FabricNodeSrc string
}

// NewConfig fetches data from the tapestry configuration file
func NewConfig() (*Config, error) {
	configFile := "tapestry.toml"
	c, err := toml.LoadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s config file: %v", configFile, err)
	}
	return &Config{
		User:          c.Get("apic.username").(string),
		Password:      c.Get("apic.password").(string),
		URL:           c.Get("apic.url").(string),
		DataSrc:       c.Get("data.src").(string),
		FabricNodeSrc: c.Get("fabricNodes.src").(string),
	}, nil
}
