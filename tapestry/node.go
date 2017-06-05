package tapestry

import (
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
