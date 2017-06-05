// Copyright © 2017 Rob Phoenix <rob@robphoenix.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/robphoenix/go-aci/aci"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get current status of ACI fabric.",
	Long:  `Get current status of ACI fabric.`,
	Run:   status,
}

func status(cmd *cobra.Command, args []string) {

	apicClient, err := aci.NewClient(Cfg.APIC.URL, Cfg.APIC.Username, Cfg.APIC.Password)
	if err != nil {
		log.Fatalf("could not create ACI client: %v", err)
	}
	err = apicClient.Login()
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}

	fmt.Printf("\nRefreshing APIC state in-memory...\n")
	fmt.Printf("\nAPIC URL: %s\n\n", apicClient.Host.Host)

	// get status of fabric nodes
	ns, err := aci.ListNodes(apicClient)
	if err != nil {
		log.Fatal(err)
	}

	// print current nodes
	fmt.Printf("Nodes\n=====\n\n")
	for _, n := range ns {
		if n.Role != "controller" {
			fmt.Printf("%s\t[ID: %s Serial: %s]\n", n.Name, n.ID, n.Serial)
		}
	}

	// get status of tenants
	ts, err := aci.ListTenants(apicClient)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nTenants\n=======\n\n")
	for _, t := range ts {
		fmt.Printf("%s\n", t.Name)
	}

}

func init() {
	RootCmd.AddCommand(statusCmd)
}
