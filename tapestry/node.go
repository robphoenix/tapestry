package tapestry

import (
	"fmt"

	"github.com/robphoenix/go-aci/aci"
)

// NodesActions ...
type NodesActions struct {
	Create []aci.Node
	Delete []aci.Node
}

// nodesStructMap builds a hash map of nodes
// indexed by Serial number
func nodesStructMap(ns []aci.Node) map[string]aci.Node {
	m := make(map[string]aci.Node, len(ns))
	for _, n := range ns {
		key := n.Serial + n.ID + n.Name
		m[key] = n
	}
	return m
}

// GetDeclaredNodes instantiates a slice of ACI Nodes from a csv data file
func GetDeclaredNodes(f string) ([]aci.Node, error) {
	data, err := CSVData(f)
	if err != nil {
		return nil, fmt.Errorf("csv data: %v", err)
	}

	var ns []aci.Node
	for _, d := range data {
		ns = append(ns, aci.Node{
			Name:   d["Name"],
			ID:     d["Node ID"],
			Serial: d["Serial"],
		})
	}
	return ns, nil
}

// DiffNodeStates determines which nodes need to be added, deleted or modified
func DiffNodeStates(desired []aci.Node, actual []aci.Node) NodesActions {
	dm := nodesStructMap(desired)
	am := nodesStructMap(actual)
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
