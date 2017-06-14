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

	"github.com/robphoenix/tapestry/tapestry"
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan changes to ACI fabric",
	Long:  `Plan changes to ACI fabric.`,
	Run:   plan,
}

func plan(cmd *cobra.Command, args []string) {

	// authenticate
	apicClient, err := tapestry.Login(Cfg.URL, Cfg.Username, Cfg.Password)
	if err != nil {
		log.Fatalf("login: %v", err)
	}

	fmt.Printf("\nRefreshing APIC state in-memory prior to plan...\n")
	fmt.Printf("\nAPIC URL: %s\n\n", apicClient.Host.Host)

	// state
	desired, err := desiredState()
	if err != nil {
		log.Fatal(err)
	}
	actual, err := actualState(apicClient)
	if err != nil {
		log.Fatal(err)
	}

	// temporary
	wantNodes := desired.Nodes()
	wantTenants := desired.Tenants()
	gotNodes := actual.Nodes()
	gotTenants := actual.Tenants()

	// actions to take
	// node
	nodeActions := tapestry.DiffNodeStates(wantNodes, gotNodes)
	fmt.Printf("Nodes\n=====\n\n")
	fmt.Printf("Plan: %d to delete, %d to create\n\n", len(nodeActions.Delete), len(nodeActions.Create))
	for _, v := range nodeActions.Delete {
		fmt.Printf("Delete -> %s [ID: %s Serial: %s]\n", v.Name, v.ID, v.Serial)
	}
	if nodeActions.Delete != nil {
		fmt.Printf("\n")
	}
	for _, v := range nodeActions.Create {
		fmt.Printf("Create -> %s [ID: %s Serial: %s]\n", v.Name, v.ID, v.Serial)
	}
	// tenant
	tenantActions := tapestry.DiffTenantStates(wantTenants, gotTenants)
	fmt.Printf("\nTenants\n=======\n\n")
	fmt.Printf("Plan: %d to delete, %d to create\n\n", len(tenantActions.Delete), len(tenantActions.Create))
	for _, v := range tenantActions.Delete {
		fmt.Printf("Delete -> %s\n", v.Name)
	}
	if tenantActions.Delete != nil {
		fmt.Printf("\n")
	}
	for _, v := range tenantActions.Create {
		fmt.Printf("Create -> %s\n", v.Name)
	}

	// summary
	fmt.Printf("\nSummary\n=======\n\n")
	fmt.Printf("Nodes: %d to delete, %d to create\n", len(nodeActions.Delete), len(nodeActions.Create))
	fmt.Printf("Tenants: %d to delete, %d to create\n", len(tenantActions.Delete), len(tenantActions.Create))
}
