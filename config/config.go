package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	APIC  `toml:"apic"`
	Nodes []Node `toml:"nodes"`
	Sites []Site `toml:"sites"`
}

func New() (Config, error) {
	var cfg Config

	viper.SetConfigName("Tapestry")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("can't read config: %v", err)
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return cfg, err
}

func NewEmpty() Config {
	return Config{
		APIC:  APIC{},
		Nodes: []Node{Node{}},
		Sites: []Site{Site{Buildings: []Building{Building{Floors: []Floor{Floor{Rooms: []Room{Room{Rows: []Row{Row{Racks: []Rack{}}}}}}}}}}},
	}
}

// TODO: config.Write()

type APIC struct {
	URL      string `toml:"url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type Node struct {
	ID     string `toml:"id"`
	Name   string `toml:"name"`
	Pod    string `toml:"pod"`
	Serial string `toml:"serial"`
	Role   string `toml:"role"`
}

type Site struct {
	Name        string     `toml:"name"`
	Description string     `toml:"description"`
	Buildings   []Building `toml:"buildings"`
}

type Building struct {
	Name        string  `toml:"name"`
	Description string  `toml:"description"`
	Floors      []Floor `toml:"floors"`
}

type Floor struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Rooms       []Room `toml:"rooms"`
}

type Room struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Rows        []Row  `toml:"rows"`
}

type Row struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Racks       []Rack `toml:"racks"`
}

type Rack struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
}
