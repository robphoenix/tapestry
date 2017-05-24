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

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan changes to ACI fabric",
	Long:  `Plan changes to ACI fabric.`,
	Run:   planChanges,
}

func planChanges(cmd *cobra.Command, args []string) {
	// create new APIC client
	apicClient, err := tapestry.NewACIClient()
	if err != nil {
		log.Fatal(err)
	}

	// fetch sources
	data := tapestry.NewSources()
	if err != nil {
		log.Fatal(err)
	}

	// read in data from fabric membership file
	fabricNodesDataFile := filepath.Join(data.FabricNodes)
	nodes, err := tapestry.NewNodes(fabricNodesDataFile)
	if err != nil {
		log.Fatal(err)
	}

	// determine desired node state
	var desiredNodeState []aci.Node
	for _, node := range nodes {
		n := aci.Node{
			Name:   node.Name,
			ID:     node.NodeID,
			Serial: node.Serial,
		}
		desiredNodeState = append(desiredNodeState, n)
	}

	// determine actual node state
	actualNodeState, err := aci.ListNodes(apicClient)
	if err != nil {
		log.Fatal(err)
	}

	// determine actions to take
	nodeActions := tapestry.DiffNodeStates(desiredNodeState, actualNodeState)

	fmt.Printf("%s\n", "Nodes to add:")
	for _, v := range nodeActions.Add {
		fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	}
	fmt.Printf("%s\n", "Nodes to delete:")
	for _, v := range nodeActions.Delete {
		fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	}

	// read in data from tenants file
	tenantsDataFile := filepath.Join(data.Tenants)
	tenants, err := tapestry.NewTenants(tenantsDataFile)
	if err != nil {
		log.Fatal(err)
	}

	// determine desired tenant state
	var desiredTenantState []aci.Tenant
	for _, tenant := range tenants {
		t := aci.Tenant{
			Name: tenant.Name,
		}
		desiredTenantState = append(desiredTenantState, t)
	}

	// determine actual tenant state
	actualTenantState, err := aci.ListTenants(apicClient)
	if err != nil {
		log.Fatal(err)
	}

	// determine actions to take
	tenantActions := tapestry.DiffTenantStates(desiredTenantState, actualTenantState)

	fmt.Printf("%s\n", "Tenants to add:")
	for _, v := range tenantActions.Add {
		fmt.Printf("%s\n", v.Name)
	}
	fmt.Printf("%s\n", "Tenants to delete:")
	for _, v := range tenantActions.Delete {
		fmt.Printf("%s\n", v.Name)
	}

}

func init() {
	RootCmd.AddCommand(planCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// planCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// planCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
