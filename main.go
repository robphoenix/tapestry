package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"github.com/pelletier/go-toml"
	"github.com/robphoenix/go-aci/aci"
)

type config struct {
	url           string
	user          string
	password      string
	dataSrc       string
	fabricNodeSrc string
}

type node struct {
	NodeID string `csv:"Node ID"`
	Name   string `csv:"Name"`
	PodID  string `csv:"Pod ID"`
	Serial string `csv:"Serial"`
	Role   string `csv:"Role"`
}

func main() {

	// cmd.Execute()

	// fetch configuration data
	conf, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}

	// create new APIC client
	apicClient, err := aci.NewClient(conf.url, conf.user, conf.password)
	if err != nil {
		log.Fatal(err)
	}

	// read in data from fabric membership file
	fabricNodesDataFile := filepath.Join(conf.dataSrc, conf.fabricNodeSrc)
	nodes, err := newNodes(fabricNodesDataFile)
	if err != nil {
		log.Fatal(err)
	}

	var ns []aci.Node
	for _, node := range nodes {
		n := aci.Node{
			Name:   node.Name,
			ID:     node.NodeID,
			Serial: node.Serial,
		}
		ns = append(ns, n)
	}

	// login
	err = apicClient.Login()
	if err != nil {
		log.Fatal(err)
	}

	// // add nodes
	// err = apicClient.AddNodes(ns)
	// if err != nil {
	//         log.Fatal(err)
	// }

	// // delete nodes
	// err = apicClient.DeleteNodes(ns)
	// if err != nil {
	//         log.Fatal(err)
	// }

	// list nodes
	n, err := apicClient.ListNodes()
	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}
	for _, node := range n {
		// fmt.Printf("node = %+v\n", node)
		fmt.Printf("Name = %+v\n", node.Name)
		// fmt.Printf("ID = %+v\n", node.ID)
		// fmt.Printf("Serial = %+v\n", node.Serial)
		// fmt.Printf("fabric status = %+v\n", node.FabricStatus)
		// fmt.Printf("node.Status = %+v\n", node.Status)
	}
}

// newConfig fetches data from the tapestry configuration file
func newConfig() (*config, error) {
	configFile := "tapestry.toml"
	c, err := toml.LoadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s config file: %v", configFile, err)
	}
	return &config{
		user:          c.Get("apic.username").(string),
		password:      c.Get("apic.password").(string),
		url:           c.Get("apic.url").(string),
		dataSrc:       c.Get("data.src").(string),
		fabricNodeSrc: c.Get("fabricNodes.src").(string),
	}, nil
}

// newNodes fetches fabric membership data from file
func newNodes(nodesFile string) ([]node, error) {
	csvFile, err := os.Open(nodesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", nodesFile, err)
	}
	defer csvFile.Close()

	var nodes []node

	err = gocsv.UnmarshalFile(csvFile, &nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal csv data: %v", err)
	}
	return nodes, nil
}
