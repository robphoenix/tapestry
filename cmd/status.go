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
	"github.com/spf13/cobra"
)

// statusCmd represents the runStatus command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get current status of ACI fabric.",
	Long:  `Get current status of ACI fabric.`,
	Run:   runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {

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

	fmt.Printf("\nCollecting current APIC state...\n")
	fmt.Printf("\nAPIC URL: %s\n\n", client.BaseURL)

	// list nodes
	nodes, err := client.FabricMembership.List(ctx)
	if err != nil {
		log.Fatalf("could not list nodes: %v", err)
	}

	// print current nodes
	fmt.Printf("Nodes\n=====\n\n")
	for _, n := range nodes {
		fmt.Printf("%s\t[ID: %s Serial: %s]\n", n.Name(), n.ID(), n.Serial())
	}
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
