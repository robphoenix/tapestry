package cmd

import (
	"fmt"

	"github.com/robphoenix/go-aci/aci"
	"github.com/robphoenix/tapestry/tapestry"
)

// DesiredState instantiates a state object, from given data,
// representing the desired state
func DesiredState() (*tapestry.State, error) {
	s := tapestry.NewState()

	// nodes
	ns, err := tapestry.GetDeclaredNodes(nodeDataFile)
	if err != nil {
		return nil, fmt.Errorf("get declared nodes: %v", err)
	}
	s.SetNodes(ns)

	// tenants
	ts, err := tapestry.GetDeclaredTenants(tenantDataFile)
	if err != nil {
		return nil, fmt.Errorf("get declared tenants: %v", err)
	}
	// set actual state
	s.SetTenants(ts)

	return s, nil
}

// ActualState instantiates a state object, from ACI API request data,
// representing the actual state
func ActualState() *tapestry.State {
	s := tapestry.NewState()

	// nodes
	ns, err := aci.ListNodes(apicClient)
	if err != nil {
		return nil, fmt.Errorf("get ACI nodes: %v", err)
	}
	// filter nodes
	ns = filterNodes(ns)
	s.SetNodes(ns)

	// actual tenant state
	ts, err := aci.ListTenants(apicClient)
	if err != nil {
		return nil, fmt.Errorf("get ACI tenants: %v", err)
	}
	// filter tenants
	ts = filterTenants(ts)
	s.SetTenants(gotTenants)

	return s, nil
}

func filterNodes(ns []aci.Node) []aci.Node {
	var fn []aci.Node
	for _, n := range ns {
		if n.Role != "controller" {
			fn = append(fn, n)
		}
	}
	return fn
}

func filterTenants(ts []aci.Tenant) []aci.Tenant {
	var ft []aci.Tenant
	for _, t := range ts {
		if t.Name != "common" && t.Name != "infra" && t.Name != "mgmt" {
			ft = append(ft, t)
		}
	}
	return ft
}
