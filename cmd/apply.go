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
	"path/filepath"

	"github.com/robphoenix/go-aci/aci"
	"github.com/robphoenix/tapestry/tapestry"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the declared state to the ACI fabric.",
	Long:  `Apply the declared state to the ACI fabric.`,
	Run:   apply,
}

func apply(cmd *cobra.Command, args []string) {

	var nCreated, nDeleted int
	var tCreated, tDeleted int

	// create new ACI client
	apicClient, err := tapestry.NewACIClient()
	if err != nil {
		log.Fatal(err)
	}

	// determine desired node state
	nf := filepath.Join(dataDir, nodesDataFile)
	wantNodes, err := tapestry.GetNodes(nf)
	if err != nil {
		log.Fatal(err)
	}

	// determine actual node state
	nodes, err := aci.ListNodes(apicClient)
	if err != nil {
		log.Fatal(err)
	}
	var gotNodes []aci.Node
	for _, n := range nodes {
		if n.Role != "controller" {
			gotNodes = append(gotNodes, n)
		}
	}

	// determine actions to take
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

	// determine desired tenant state
	tf := filepath.Join(dataDir, tenantsDataFile)
	wantTenants, err := tapestry.GetTenants(tf)
	if err != nil {
		log.Fatal(err)
	}

	// determine actual tenant state
	tenants, err := aci.ListTenants(apicClient)
	if err != nil {
		log.Fatal(err)
	}
	var gotTenants []aci.Tenant
	for _, t := range tenants {
		if t.Name != "common" && t.Name != "infra" && t.Name != "mgmt" {
			gotTenants = append(gotTenants, t)
		}
	}

	// determine actions to take
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

	fmt.Printf("\nSummary\n=======\n\n")
	fmt.Printf("Nodes: %d deleted, %d created\n", nDeleted, nCreated)
	fmt.Printf("Tenants: %d deleted, %d created\n", tDeleted, tCreated)
}

func init() {
	RootCmd.AddCommand(applyCmd)
}
