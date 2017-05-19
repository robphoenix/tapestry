package tapestry

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/robphoenix/go-aci/aci"
)

// Node ...
type Node struct {
	NodeID string `csv:"Node ID"`
	Name   string `csv:"Name"`
	PodID  string `csv:"Pod ID"`
	Serial string `csv:"Serial"`
	Role   string `csv:"Role"`
}

// NodesActions ...
type NodesActions struct {
	Add    []aci.Node
	Delete []aci.Node
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

// DiffNodeStates determines which nodes need to be added, deleted or modified
func DiffNodeStates(desired []aci.Node, actual []aci.Node) NodesActions {
	dm := structMap(desired)
	am := structMap(actual)
	var na NodesActions

	// add
	for k, dv := range dm {
		_, ok := am[k]
		if !ok {
			na.Add = append(na.Add, dv)
		}
	}
	// delete
	for k, av := range am {
		_, ok := dm[k]
		if !ok {
			na.Delete = append(na.Delete, av)
		}
	}
	return na
}

// NewNodes fetches fabric membership data from file
func NewNodes(nodesFile string) ([]Node, error) {
	csvFile, err := os.Open(nodesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", nodesFile, err)
	}
	defer csvFile.Close()

	var nodes []Node

	err = gocsv.UnmarshalFile(csvFile, &nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal csv data: %v", err)
	}
	return nodes, nil
}
