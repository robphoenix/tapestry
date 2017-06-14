package tapestry

import (
	"fmt"

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

// GetDeclaredTenants fetches tenant data from file
func GetDeclaredTenants(f string) ([]aci.Tenant, error) {
	data, err := CSVData(f)
	if err != nil {
		return nil, fmt.Errorf("csv data: %v", err)
	}

	var ts []aci.Tenant
	for _, d := range data {
		ts = append(ts, aci.Tenant{Name: d["Name"]})
	}
	return ts, nil
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
