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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	nodeDataFile   = "data/fabric_membership.csv"
	tenantDataFile = "data/tenant.csv"
)

var (
	// Cfg holds APIC configuration details
	Cfg config
	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "tapestry",
		Short: "Weave a Cisco ACI fabric",
		Long:  `Tapestry is a CLI tool for declaring and deploying your Cisco ACI fabric.`,
	}
)

type config struct {
	APIC
}

// APIC holds all the APIC configuration details
type APIC struct {
	URL      string
	Username string
	Password string
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	// add child commands to root command
	RootCmd.AddCommand(statusCmd)
	RootCmd.AddCommand(planCmd)
	RootCmd.AddCommand(applyCmd)
	RootCmd.AddCommand(destroyCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file
func initConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // look for config in the working directory

	// read in config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	// unmarshal config data
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
