package tapestry

import (
	"github.com/robphoenix/go-aci/aci"
)

// TenantsActions ...
type TenantsActions struct {
	Create []aci.Tenant
	Delete []aci.Tenant
}

// tenantsStructMap builds a hash map of tenants
// indexed by name
func tenantsStructMap(ts []aci.Tenant) map[string]aci.Tenant {
	m := make(map[string]aci.Tenant, len(ts))
	for _, t := range ts {
		m[t.Name] = t
	}
	return m
}

// DiffTenantStates determines which tenants need to be added, deleted or modified
func DiffTenantStates(desired, actual []aci.Tenant) TenantsActions {
	dm := tenantsStructMap(desired)
	am := tenantsStructMap(actual)
	var ta TenantsActions

	// add
	for k, dv := range dm {
		_, ok := am[k]
		if !ok {
			ta.Create = append(ta.Create, dv)
		}
	}
	// delete
	for k, av := range am {
		_, ok := dm[k]
		if !ok {
			ta.Delete = append(ta.Delete, av)
		}
	}
	return ta
}
