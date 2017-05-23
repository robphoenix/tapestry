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

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get current status of ACI fabric.",
	Long:  `Get current status of ACI fabric.`,
	Run:   aciStatus,
}

func aciStatus(cmd *cobra.Command, args []string) {

	// get status of fabric nodes
	nodeStatus()
	// get status of tenants
	tenantStatus()

}

func nodeStatus() {
	c, err := tapestry.NewACIClient()
	if err != nil {
		log.Fatal(err)
	}

	ns, err := aci.ListNodes(c)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", "Current Nodes:")
	for _, n := range ns {
		fmt.Printf("%s\t%s\t%s\n", n.Name, n.ID, n.Serial)
	}
}

func tenantStatus() {
	c, err := tapestry.NewACIClient()
	if err != nil {
		log.Fatal(err)
	}

	ts, err := aci.ListTenants(c)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", "Current Tenants:")
	for _, t := range ts {
		fmt.Printf("%s\n", t.Name)
	}
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
