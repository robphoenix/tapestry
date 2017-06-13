package tapestry

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/robphoenix/go-aci/aci"
)

const (
	tenantDataFile = "data/tenant.csv"
)

// Tenant ...
type Tenant struct {
	Name string `csv:"Name"`
}

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
func GetDeclaredTenants() ([]aci.Tenant, error) {
	csvFile, err := os.Open(tenantDataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", tenantDataFile, err)
	}
	defer csvFile.Close()

	var ts []Tenant

	err = gocsv.UnmarshalFile(csvFile, &ts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal csv data: %v", err)
	}

	var ats []aci.Tenant
	for _, t := range ts {
		ats = append(ats, aci.Tenant{Name: t.Name})
	}
	return ats, nil
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
