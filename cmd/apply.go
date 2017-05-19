// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fetch configuration data
		conf, err := tapestry.NewConfig()
		if err != nil {
			log.Fatal(err)
		}

		// create new APIC client
		apicClient, err := aci.NewClient(conf.URL, conf.User, conf.Password)
		if err != nil {
			log.Fatal(err)
		}

		// read in data from fabric membership file
		fabricNodesDataFile := filepath.Join(conf.DataSrc, conf.FabricNodeSrc)
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

		// login
		err = apicClient.Login()
		if err != nil {
			log.Fatal(err)
		}

		// determine actual node state
		actualNodeState, err := apicClient.ListNodes()
		if err != nil {
			log.Fatal(err)
		}

		// determine actions to take
		na := tapestry.DiffNodeStates(desiredNodeState, actualNodeState)

		// delete nodes
		err = apicClient.DeleteNodes(na.Delete)
		if err != nil {
			log.Fatal(err)
		}
		if err == nil {
			fmt.Printf("%s\n", "Nodes deleted:")
			for _, v := range na.Delete {
				fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
			}
		}

		// add nodes
		err = apicClient.AddNodes(na.Add)
		if err != nil {
			log.Fatal(err)
		}
		if err == nil {
			fmt.Printf("%s\n", "Nodes added:")
			for _, v := range na.Add {
				fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
