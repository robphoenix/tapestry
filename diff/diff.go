package diff

import (
	set "github.com/deckarep/golang-set"
	"github.com/robphoenix/go-aci/aci"
)

// CompareNodes compares two slices of ACI nodes, identifying
// what changes are to be made to bring ACI into the desired
// state.
func CompareNodes(a, b []*aci.Node) (add, delete []*aci.Node) {
	var x []interface{}
	var y []interface{}

	for _, v := range a {
		x = append(x, v)
	}
	for _, v := range b {
		y = append(y, v)
	}

	xx := set.NewSetFromSlice(x)
	yy := set.NewSetFromSlice(y)
	aa := xx.Difference(yy)
	bb := yy.Difference(xx)

	for elem := range aa.Iterator().C {
		add = append(add, elem.(*aci.Node))
	}
	for elem := range bb.Iterator().C {
		delete = append(delete, elem.(*aci.Node))
	}

	return add, delete
}
