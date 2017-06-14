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
	"github.com/robphoenix/tapestry/tapestry"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the declared state to the ACI fabric.",
	Long:  `Apply the declared state to the ACI fabric.`,
	Run:   apply,
}

func apply(cmd *cobra.Command, args []string) {

	// declare counters
	var nCreated, nDeleted, tCreated, tDeleted int

	// authenticate
	apicClient, err := tapestry.Login(Cfg.URL, Cfg.Username, Cfg.Password)
	if err != nil {
		log.Fatalf("login: %v", err)
	}

	// state
	desired, err := desiredState()
	if err != nil {
		log.Fatal(err)
	}
	actual, err := actualState(apicClient)
	if err != nil {
		log.Fatal(err)
	}

	wantNodes := desired.Nodes()
	wantTenants := desired.Tenants()
	gotNodes := actual.Nodes()
	gotTenants := actual.Tenants()

	// actions to take
	nodeActions := tapestry.DiffNodeStates(wantNodes, gotNodes)
	// delete nodes
	// do this first as we can't modify nodes that already exist
	// so we have to delete and then re-add them
	if nodeActions.Delete != nil {
		err = aci.DeleteNodes(apicClient, nodeActions.Delete)
		if err != nil {
			log.Fatal(err)
		}
		if err == nil {
			nDeleted = len(nodeActions.Delete)
			fmt.Printf("Deleting Nodes...\n\n")
			for _, v := range nodeActions.Delete {
				fmt.Printf("%s [ID: %s Serial: %s]\n", v.Name, v.ID, v.Serial)
			}
		}
	}
	// create nodes
	if nodeActions.Create != nil {
		err = aci.CreateNodes(apicClient, nodeActions.Create)
		if err != nil {
			log.Fatal(err)
		}
		if err == nil {
			nCreated = len(nodeActions.Create)
			fmt.Printf("Creating Nodes...\n\n")
			for _, v := range nodeActions.Create {
				fmt.Printf("%s [ID: %s Serial: %s]\n", v.Name, v.ID, v.Serial)
			}
		}
	}

	// actions to take
	tenantActions := tapestry.DiffTenantStates(wantTenants, gotTenants)
	// delete tenants
	// do this first as we can't modify nodes that already exist
	// so we have to delete and then re-add them
	if tenantActions.Delete != nil {
		fmt.Printf("\nDeleting Tenants...\n\n")
		for _, v := range tenantActions.Delete {
			err = aci.DeleteTenant(apicClient, v)
			if err != nil {
				log.Fatal(err)
			}
			if err == nil {
				tDeleted++
				fmt.Printf("%s\n", v.Name)
			}
		}
	}
	// create tenants
	if tenantActions.Create != nil {
		fmt.Printf("\nCreating Tenants...\n\n")
		for _, v := range tenantActions.Create {
			err = aci.CreateTenant(apicClient, v)
			if err != nil {
				log.Fatal(err)
			}
			if err == nil {
				tCreated++
				fmt.Printf("%s\n", v.Name)
			}
		}
	}

	// summary
	fmt.Printf("\nSummary\n=======\n\n")
	fmt.Printf("Nodes: %d deleted, %d created\n", nDeleted, nCreated)
	fmt.Printf("Tenants: %d deleted, %d created\n", tDeleted, tCreated)
}
