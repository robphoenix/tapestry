package tapestry

import (
	"fmt"

	"github.com/robphoenix/go-aci/aci"
	"github.com/spf13/viper"
)

// NewACIClient creates a new authenticated
// ACI client from given configuration
func NewACIClient() (*aci.Client, error) {
	// fetch configuration data
	apicURL := viper.GetString("apic.url")
	apicUser := viper.GetString("apic.username")
	apicPwd := viper.GetString("apic.password")

	// create new APIC client
	apicClient, err := aci.NewClient(apicURL, apicUser, apicPwd)
	if err != nil {
		return nil, fmt.Errorf("could not create ACI client: %v", err)
	}

	// login
	err = apicClient.Login()
	if err != nil {
		return nil, fmt.Errorf("could not login: %v", err)
	}

	return apicClient, nil
}
