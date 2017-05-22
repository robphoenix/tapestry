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

	"github.com/robphoenix/go-aci/aci"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// fetch configuration data
		apicURL := viper.GetString("apic.url")
		apicUser := viper.GetString("apic.username")
		apicPwd := viper.GetString("apic.password")

		// create new APIC client
		apicClient, err := aci.NewClient(apicURL, apicUser, apicPwd)
		if err != nil {
			log.Fatal(err)
		}

		// login
		err = apicClient.Login()
		if err != nil {
			log.Fatal(err)
		}

		// determine node state
		nodeState, err := apicClient.ListNodes()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", "Current Nodes:")
		for _, v := range nodeState {
			fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
