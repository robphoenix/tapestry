package tapestry

import "github.com/robphoenix/go-aci/aci"

// State describes an ACI state, whether desired or actual
type State struct {
	nodes   []aci.Node
	tenants []aci.Tenant
}

// NewState instantiates a new State object
func NewState() *State {
	return &State{}
}

// Nodes lists the fabric membership nodes the state has
func (s *State) Nodes() []aci.Node {
	return s.nodes
}

// SetNodes sets the states fabric membership nodes
func (s *State) SetNodes(ns []aci.Node) {
	s.nodes = ns
}

// AddNode adds an ACI node to the state
func (s *State) AddNode(n aci.Node) {
	s.nodes = append(s.nodes, n)
}

// DelNode deletes the given ACI node from the state
func (s *State) DelNode(n aci.Node) {
	for i, node := range s.nodes {
		if node == n {
			s.nodes = append(s.nodes[:i], s.nodes[i+1:]...)
			break
		}
	}
}

// DelAllNodes deletes all ACI nodes from the state
func (s *State) DelAllNodes() {
	s.nodes = nil
}

// HasNode returns whether a node is present in the state or not
func (s *State) HasNode(n aci.Node) bool {
	for _, node := range s.nodes {
		if node == n {
			return true
		}
	}
	return false
}

// Tenants lists the tenants present in the state
func (s *State) Tenants() []aci.Tenant {
	return s.tenants
}

// SetTenants sets the state's tenants
func (s *State) SetTenants(ts []aci.Tenant) {
	s.tenants = ts
}

// AddTenant adds an ACI node to the state
func (s *State) AddTenant(t aci.Tenant) {
	s.tenants = append(s.tenants, t)
}

// DelTenant deletes the given ACI node from the state
func (s *State) DelTenant(t aci.Tenant) {
	for i, tenant := range s.tenants {
		if tenant == t {
			s.tenants = append(s.tenants[:i], s.tenants[i+1:]...)
			break
		}
	}
}

// DelAllTenant deletes all ACI tenants from the state
func (s *State) DelAllTenant() {
	s.tenants = nil
}

// HasTenant returns whether a tenant is present in the state or not
func (s *State) HasTenant(t aci.Tenant) bool {
	for _, tenant := range s.tenants {
		if tenant == t {
			return true
		}
	}
	return false
}
