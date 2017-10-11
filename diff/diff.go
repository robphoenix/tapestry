package diff

import (
	"github.com/robphoenix/go-aci/aci"
)

// CompareNodes compares two slices of ACI nodes, identifying
// what changes are to be made to bring ACI into the desired
// state.
func CompareNodes(a, b []*aci.Node) (add, delete []*aci.Node) {
	for _, aa := range a {
		if !containsNode(aa, b) {
			add = append(add, aa)
		}
	}
	for _, bb := range b {
		if !containsNode(bb, a) {
			delete = append(delete, bb)
		}
	}

	return add, delete
}

func containsNode(n *aci.Node, nodes []*aci.Node) bool {
	for _, node := range nodes {
		if node.Equal(n) {
			return true
		}
	}
	return false
}
