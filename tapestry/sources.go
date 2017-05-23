package tapestry

import (
	"fmt"
	"path/filepath"

	"github.com/robphoenix/go-aci/aci"
	"github.com/spf13/viper"
)

// Sources stores the absolute paths to data sources
type Sources struct {
	DataDir     string
	FabricNodes string
	Tenants     []aci.Tenant
}

// NewSources instantiates a new Sources struct
func NewSources() (Sources, error) {
	dd := viper.GetString("data.src")
	s := Sources{
		DataDir:     dd,
		FabricNodes: filepath.Join(dd, viper.GetString("fabricnodes.src")),
	}
	err := viper.UnmarshalKey("tenants", &s.Tenants)
	if err != nil {
		return s, err
	}
	fmt.Printf("s = %+v\n", s)
	return s, nil
}
