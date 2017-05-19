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

type nodesActions struct {
	add    []aci.Node
	delete []aci.Node
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

	// determine desired node state
	var desiredNodeState []aci.Node
	for _, node := range nodes {
		n := aci.Node{
			Name:   node.Name,
			ID:     node.NodeID,
			Serial: node.Serial,
		}
		desiredNodeState = append(desiredNodeState, n)
	}

	// login
	err = apicClient.Login()
	if err != nil {
		log.Fatal(err)
	}

	// determine actual node state
	actualNodeState, err := apicClient.ListNodes()
	if err != nil {
		log.Fatal(err)
	}

	// determine actions to take
	na := diffNodeStates(desiredNodeState, actualNodeState)

	fmt.Printf("%s\n", "Nodes to add:")
	for _, v := range na.add {
		fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	}
	fmt.Printf("%s\n", "Nodes to delete:")
	for _, v := range na.delete {
		fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	}

	// delete nodes
	err = apicClient.DeleteNodes(na.delete)
	if err != nil {
		log.Fatal(err)
	}
	if err == nil {
		fmt.Printf("%s\n", "Nodes deleted:")
		for _, v := range na.delete {
			fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
		}
	}

	// add nodes
	err = apicClient.AddNodes(na.add)
	if err != nil {
		log.Fatal(err)
	}
	if err == nil {
		fmt.Printf("%s\n", "Nodes added:")
		for _, v := range na.add {
			fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
		}
	}

	// determine new node state
	newNodeState, err := apicClient.ListNodes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", "Current Nodes:")
	for _, v := range newNodeState {
		fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
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

// diffNodeStates determines which nodes need to be added, deleted or modified
func diffNodeStates(desired []aci.Node, actual []aci.Node) nodesActions {
	dm := structMap(desired)
	am := structMap(actual)
	var na nodesActions

	// add
	for k, dv := range dm {
		_, ok := am[k]
		if !ok {
			na.add = append(na.add, dv)
		}
	}
	// delete
	for k, av := range am {
		_, ok := dm[k]
		if !ok {
			na.delete = append(na.delete, av)
		}
	}
	return na
}

// structMap builds a hash map of nodes
// indexed by Serial number
func structMap(ns []aci.Node) map[string]aci.Node {
	m := make(map[string]aci.Node, len(ns))
	for _, n := range ns {
		key := n.Serial + n.ID + n.Name
		m[key] = n
	}
	return m
}
