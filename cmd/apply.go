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

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the declared state to the ACI fabric.",
	Long:  `Apply the declared state to the ACI fabric.`,
	Run: func(cmd *cobra.Command, args []string) {
		applyState()
	},
}

func applyState() {

	// create new APIC client
	apicClient, err := tapestry.NewACIClient()
	if err != nil {
		log.Fatal(err)
	}

	// read in data from fabric membership file
	nodes, err := tapestry.NewNodes(tapestry.NewSources().FabricNodes)
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
	actions := tapestry.DiffNodeStates(desiredNodeState, actualNodeState)

	// delete nodes
	err = aci.DeleteNodes(apicClient, actions.Delete)
	if err != nil {
		log.Fatal(err)
	}
	if err == nil {
		fmt.Printf("%s\n", "Nodes deleted:")
		for _, v := range actions.Delete {
			fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
		}
	}

	// add nodes
	err = aci.AddNodes(apicClient, actions.Add)
	if err != nil {
		log.Fatal(err)
	}
	if err == nil {
		fmt.Printf("%s\n", "Nodes added:")
		for _, v := range actions.Add {
			fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
		}
	}

}

func init() {
	RootCmd.AddCommand(applyCmd)
}
