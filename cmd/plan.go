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

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan changes to ACI fabric",
	Long:  `Plan changes to ACI fabric.`,
	Run:   plan,
}

func plan(cmd *cobra.Command, args []string) {

	apicClient, err := aci.NewClient(Cfg.APIC.URL, Cfg.APIC.Username, Cfg.APIC.Password)
	if err != nil {
		log.Fatalf("could not create ACI client: %v", err)
	}
	err = apicClient.Login()
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}

	fmt.Printf("\nRefreshing APIC state in-memory prior to plan...\n")
	fmt.Printf("\nAPIC URL: %s\n\n", apicClient.Host.Host)

	// actual node state
	aciNodes, err := aci.ListNodes(apicClient)
	if err != nil {
		log.Fatal(err)
	}
	// filter out controllers
	var gotNodes []aci.Node
	for _, n := range aciNodes {
		if n.Role != "controller" {
			gotNodes = append(gotNodes, n)
		}
	}

	// desired node state
	nodes := Cfg.Nodes
	var wantNodes []aci.Node
	for _, n := range nodes {
		wn := aci.Node{
			Name:   n.Name,
			ID:     n.ID,
			Serial: n.Serial,
		}
		wantNodes = append(wantNodes, wn)
	}

	// actions to take
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

	// desired tenant state
	tenants := Cfg.Tenants
	var wantTenants []aci.Tenant
	for _, t := range tenants {
		wt := aci.Tenant{
			Name: t.Name,
		}
		wantTenants = append(wantTenants, wt)
	}

	// determine actual tenant state
	aciTenants, err := aci.ListTenants(apicClient)
	if err != nil {
		log.Fatal(err)
	}
	var gotTenants []aci.Tenant
	for _, t := range aciTenants {
		if t.Name != "common" && t.Name != "infra" && t.Name != "mgmt" {
			gotTenants = append(gotTenants, t)
		}
	}

	// determine actions to take
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

	fmt.Printf("\nSummary\n=======\n\n")
	fmt.Printf("Nodes: %d to delete, %d to create\n", len(nodeActions.Delete), len(nodeActions.Create))
	fmt.Printf("Tenants: %d to delete, %d to create\n", len(tenantActions.Delete), len(tenantActions.Create))
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
