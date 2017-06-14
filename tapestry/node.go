package tapestry

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/robphoenix/go-aci/aci"
)

const (
	dataFile = "data/fabric_membership.csv"
)

// Node ...
type Node struct {
	Name   string `csv:"Name"`
	NodeID string `csv:"Node ID"`
	PodID  string `csv:"Pod ID"`
	Role   string `csv:"Role"`
	Serial string `csv:"Serial"`
}

// NodesActions ...
type NodesActions struct {
	Create []aci.Node
	Delete []aci.Node
}

// GetDeclaredNodes instantiates a new Nodes struct from a csv data file
func GetDeclaredNodes() ([]aci.Node, error) {
	csvFile, err := os.Open(dataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", dataFile, err)
	}
	defer csvFile.Close()

	var ns []Node

	err = gocsv.UnmarshalFile(csvFile, &ns)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal csv data: %v", err)
	}

	var ans []aci.Node
	for _, n := range ns {
		ans = append(ans, aci.Node{
			Name:   n.Name,
			ID:     n.NodeID,
			Serial: n.Serial,
		})
	}
	return ans, nil
}

// DiffNodeStates determines which nodes need to be added, deleted or modified
func DiffNodeStates(desired []aci.Node, actual []aci.Node) NodesActions {
	dm := make(map[string]aci.Node, len(desired))
	for _, d := range desired {
		dm[d.Key()] = d.Value()
	}
	am := make(map[string]aci.Node, len(desired))
	for _, a := range actual {
		am[a.Key()] = a.Value()
	}
	var na NodesActions

	// add
	for k, dv := range dm {
		_, ok := am[k]
		if !ok {
			na.Create = append(na.Create, dv)
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
