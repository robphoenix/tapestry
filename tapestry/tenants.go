package tapestry

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/robphoenix/go-aci/aci"
)

// Tenant ...
type Tenant struct {
	Name string `csv:"Name"`
}

// TenantsActions ...
type TenantsActions struct {
	Add    []aci.Tenant
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
func DiffTenantStates(desired []aci.Tenant, actual []aci.Tenant) TenantsActions {
	dm := tenantsStructMap(desired)
	am := tenantsStructMap(actual)
	var ta TenantsActions

	// add
	for k, dv := range dm {
		_, ok := am[k]
		if !ok {
			ta.Add = append(ta.Add, dv)
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

// NewTenants fetches tenant data from file
func NewTenants(tenantsFile string) ([]Tenant, error) {
	csvFile, err := os.Open(tenantsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", tenantsFile, err)
	}
	defer csvFile.Close()

	var tenants []Tenant

	err = gocsv.UnmarshalFile(csvFile, &tenants)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal csv data: %v", err)
	}
	return tenants, nil
}
