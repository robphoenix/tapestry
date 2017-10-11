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
	"github.com/robphoenix/tapestry/diff"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the declared state to the ACI fabric.",
	Long:  `Apply the declared state to the ACI fabric.`,
	Run:   runApply,
}

func runApply(cmd *cobra.Command, args []string) {

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

	add, delete := diff.CompareNodes(wantNodes, gotNodes)

	if len(add) == 0 && len(delete) == 0 {
		fmt.Println("No changes to be made")
	}
	if len(delete) > 0 {
		fmt.Printf("Deleting Nodes\n")
		fmt.Printf("==============\n\n")
		for _, v := range delete {
			fmt.Printf("%+v\n", v)
			v.SetDeleted()
		}
		fmt.Println("")
	}
	if len(add) > 0 {
		fmt.Printf("Adding Nodes\n")
		fmt.Printf("============\n\n")
		for _, v := range add {
			fmt.Printf("%+v\n", v)
			v.SetCreated()
		}
	}

	// update ACI Fabric Membership
	// TODO: will fail if any nodes need to be decommissioned first.
	// we must delete nodes first, otherwise we vould have duplicate node id's
	_, err = client.FabricMembership.Update(ctx, delete...)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = client.FabricMembership.Update(ctx, add...)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func init() {
	RootCmd.AddCommand(applyCmd)
}
