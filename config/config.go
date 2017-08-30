package config

// Config ...
type Config struct {
	APIC
	Nodes []Node
	Sites []Site
}

type APIC struct {
	URL      string
	Username string
	Password string
}

type Node struct {
	ID     string
	Name   string
	Pod    string
	Serial string
	Role   string
}

type Site struct {
	Name        string
	Description string
	Buildings   []Building
}

type Building struct {
	Name        string
	Description string
	Floors      []Floor
}

type Floor struct {
	Name        string
	Description string
	Rooms       []Room
}

type Room struct {
	Name        string
	Description string
	Rows        []Row
}

type Row struct {
	Name        string
	Description string
	Racks       []Rack
}

type Rack struct {
	Name        string
	Description string
}
