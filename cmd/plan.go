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
	"context"
	"fmt"
	"log"

	"github.com/robphoenix/go-aci/aci"
	"github.com/robphoenix/tapestry/config"
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan changes to ACI fabric",
	Long:  `Plan changes to ACI fabric.`,
	Run:   runPlan,
}

func runPlan(cmd *cobra.Command, args []string) {

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	client, err := aci.NewClient(aci.Config{
		Host:     cfg.URL,
		Username: cfg.Username,
		Password: cfg.Password,
	})
	if err != nil {
		log.Fatalf("could not create ACI client: %v", err)
	}

	ctx := context.Background()

	err = client.Login(ctx)
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}

	fmt.Printf("\nRefreshing APIC state in-memory...\n")
	fmt.Printf("\nAPIC URL: %s\n\n", client.BaseURL)

	// desired node state
	var wantNodes []*aci.Node
	for _, node := range cfg.Nodes {
		n, err := client.FabricMembership.NewNode(
			node.Name,
			node.ID,
			node.Pod,
			node.Serial,
			node.Role,
		)
		if err != nil {
			log.Fatalf("invalid node: %v", err)
		}
		wantNodes = append(wantNodes, n)
	}

	fmt.Printf("Desired Nodes\n=====\n\n")
	for _, n := range wantNodes {
		fmt.Printf("%s\t[ID: %s Serial: %s]\n", n.Name(), n.ID(), n.Serial())
	}
	fmt.Printf("\n")

	// actual node state
	nodes, err := client.FabricMembership.List(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	var gotNodes []*aci.Node
	for _, n := range nodes {
		if n.Role() != "controller" {
			gotNodes = append(gotNodes, n)
		}
	}

	// print actual nodes
	fmt.Printf("Actual Nodes\n=====\n\n")
	for _, n := range gotNodes {
		fmt.Printf("%s\t[ID: %s Serial: %s]\n", n.Name(), n.ID(), n.Serial())
	}
	fmt.Printf("\n")

	//TODO diff desired & actual

	//         // summary
	//         fmt.Printf("\nSummary\n=======\n\n")
	//         fmt.Printf("Nodes: %d to delete, %d to create\n", len(nodeActions.Delete), len(nodeActions.Create))
	//         fmt.Printf("Tenants: %d to delete, %d to create\n", len(tenantActions.Delete), len(tenantActions.Create))
}

func init() {
	RootCmd.AddCommand(planCmd)
}
