package tapestry

import (
	"fmt"

	"github.com/robphoenix/go-aci/aci"
)

// Login authenticates with the ACI server
func Login(u, n, p string) (*aci.Client, error) {
	c, err := aci.NewClient(u, n, p)
	if err != nil {
		return nil, fmt.Errorf("could not create ACI client: %v", err)
	}
	err = c.Login()
	if err != nil {
		return nil, fmt.Errorf("could not login: %v", err)
	}
	return c, nil
}
