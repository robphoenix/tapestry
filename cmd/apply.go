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
	Run:   applyState,
}

func applyState(cmd *cobra.Command, args []string) {

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
	gotNodes, err := aci.ListNodes(apicClient)
	if err != nil {
		log.Fatal(err)
	}

	// determine actions to take
	nodeActions := tapestry.DiffNodeStates(wantNodes, gotNodes)

	fmt.Printf("%s\n", "Nodes to add:")
	for _, v := range nodeActions.Add {
		fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	}
	fmt.Printf("%s\n", "Nodes to delete:")
	for _, v := range nodeActions.Delete {
		fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	}

	// delete nodes
	err = aci.DeleteNodes(apicClient, nodeActions.Delete)
	if err != nil {
		log.Fatal(err)
	}
	if err == nil {
		fmt.Printf("%s\n", "Nodes deleted:")
		for _, v := range nodeActions.Delete {
			fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
		}
	}

	// add nodes
	err = aci.AddNodes(apicClient, nodeActions.Add)
	if err != nil {
		log.Fatal(err)
	}
	if err == nil {
		fmt.Printf("%s\n", "Nodes added:")
		for _, v := range nodeActions.Add {
			fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
		}
	}

	// determine desired tenant state
	tf := filepath.Join(dataDir, tenantsDataFile)
	wantTenants, err := tapestry.GetTenants(tf)
	if err != nil {
		log.Fatal(err)
	}

	// determine actual tenant state
	gotTenants, err := aci.ListTenants(apicClient)
	if err != nil {
		log.Fatal(err)
	}

	// determine actions to take
	tenantActions := tapestry.DiffTenantStates(wantTenants, gotTenants)

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
	RootCmd.AddCommand(applyCmd)
}
