package cmd

import (
	"fmt"

	"github.com/robphoenix/go-aci/aci"
)

func login() (*aci.Client, error) {
	c, err := aci.NewClient(Cfg.URL, Cfg.Username, Cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("could not create ACI client: %v", err)
	}
	err = c.Login()
	if err != nil {
		return nil, fmt.Errorf("could not login: %v", err)
	}
	return c, nil
}
