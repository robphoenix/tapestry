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
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/robphoenix/go-aci/aci"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		destroy()
	},
}

func destroy() {

	prompt := "Do you really want to destroy?\n\n" +
		"Tapestry will delete all your APIC configuration.\n" +
		"There is no undo. Only 'yes' will be accepted to confirm.\n\n" +
		"Enter value: "
	fmt.Printf(prompt)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %#v\n", err.Error())
		os.Exit(1)
	}

	if strings.TrimSpace(response) != "yes" {
		fmt.Printf("Destroy Cancelled.\n")
		os.Exit(1)
	}

	var nDeleted, tDeleted int

	apicClient, err := aci.NewClient(Cfg.APIC.URL, Cfg.APIC.Username, Cfg.APIC.Password)
	if err != nil {
		log.Fatalf("could not create ACI client: %v", err)
	}
	err = apicClient.Login()
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}

	// get node state
	aciNodes, err := aci.ListNodes(apicClient)
	if err != nil {
		log.Fatal(err)
	}
	var nodes []aci.Node
	for _, n := range aciNodes {
		if n.Role != "controller" {
			nodes = append(nodes, n)
		}
	}

	// delete nodes
	if nodes != nil {
		err = aci.DeleteNodes(apicClient, nodes)
		if err != nil {
			log.Fatal(err)
		}
		if err == nil {
			nDeleted = len(nodes)
			fmt.Printf("Deleting Nodes...\n\n")
			for _, v := range nodes {
				fmt.Printf("%s [ID: %s Serial: %s]\n", v.Name, v.ID, v.Serial)
			}
		}
	}

	// get tenant state
	aciTenants, err := aci.ListTenants(apicClient)
	if err != nil {
		log.Fatal(err)
	}
	var tenants []aci.Tenant
	for _, t := range aciTenants {
		if t.Name != "common" && t.Name != "infra" && t.Name != "mgmt" {
			tenants = append(tenants, t)
		}
	}

	// delete tenants
	if tenants != nil {
		fmt.Printf("\nDeleting Tenants...\n\n")
		for _, v := range tenants {
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

	fmt.Printf("\nSummary\n=======\n\n")
	fmt.Printf("Nodes: %d deleted\n", nDeleted)
	fmt.Printf("Tenants: %d deleted\n", tDeleted)

}

func init() {
	RootCmd.AddCommand(destroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// destroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
