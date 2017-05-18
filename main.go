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
	modify []aci.Node
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

	na := diffNodeStates(desiredNodeState, actualNodeState)
	fmt.Printf("na.add = %+v\n", na.add)
	fmt.Printf("na.delete = %+v\n", na.delete)
	fmt.Printf("na.modify = %+v\n", na.modify)

	// // add nodes
	// err = apicClient.AddNodes(desiredNodeState)
	// if err != nil {
	//         log.Fatal(err)
	// }

	// // delete nodes
	// err = apicClient.DeleteNodes(desiredNodeState)
	// if err != nil {
	//         log.Fatal(err)
	// }

	// list nodes
	for _, node := range actualNodeState {
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

// diffNodeStates determines which nodes need to be added, deleted or modified
func diffNodeStates(desired []aci.Node, actual []aci.Node) nodesActions {
	dm := structMap(desired)
	am := structMap(actual)
	var na nodesActions

	// add
	for k, v := range dm {
		_, ok := am[k]
		if !ok {
			na.add = append(na.add, v)
		}
	}
	// delete
	for k, v := range am {
		_, ok := dm[k]
		if !ok {
			na.delete = append(na.delete, v)
		}
	}
	// modify
	for k, dv := range dm {
		av, ok := am[k]
		if ok {
			if dv.Name != av.Name || dv.ID != av.ID {
				na.modify = append(na.modify, dv)
			}
		}
	}
	return na
}

// structMap builds a hash map of nodes
// indexed by Serial number
func structMap(ns []aci.Node) map[string]aci.Node {
	m := make(map[string]aci.Node, len(ns))
	for _, n := range ns {
		m[n.Serial] = n
	}
	return m
}
