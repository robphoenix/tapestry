package tapestry

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// Sources stores the absolute paths to data sources
type Sources struct {
	DataDir     string
	FabricNodes string
}

// NewSources instantiates a new Sources struct
func NewSources() Sources {
	dd := viper.GetString("data.src")
	return Sources{
		DataDir:     dd,
		FabricNodes: filepath.Join(dd, viper.GetString("fabricnodes.src")),
	}
}
