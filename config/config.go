package config

// Config ...
type Config struct {
	APIC  `toml:"apic"`
	Nodes []Node `toml:"nodes"`
	Sites []Site `toml:"sites"`
}

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
