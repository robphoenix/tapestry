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
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new",
	Short:   "Create a new Tapestry project.",
	Long:    `New creates a Tapestry project with the given name.`,
	Example: "tapestry new <project>",
	Args:    cobra.ExactArgs(1),
	Run:     runNew,
}

func runNew(cmd *cobra.Command, args []string) {
	// TODO validate directory names
	dir := args[0]

	if info, _ := os.Stat(dir); info != nil {
		fmt.Printf("directory already exists: %s\n", dir)
		return
	}

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filepath.Join(dir, "Tapestry.toml"))
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("config").Parse(content)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(f, args[0])

	fmt.Printf("new Tapestry project created: %s\n", dir)
}

const content = `
# {{.}}

# ACI APIC
[apic]
url = "example-apic.com"
username = "admin"
password = "password"

# Fabric Membership
[[nodes]]
id = 101
name = "leaf-901"
pod = 1
serial = "ABCDEF12345"
role = "leaf"

# Geolocation
[[sites]]
name = "site-01"
description = "Site 01"

[[sites.buildings]]
name = "building-01"
description = "Building 01"

[[sites.buildings.floors]]
name = "floor-01"
description = "Floor 01"

[[sites.buildings.floors.rooms]]
name = "room-01"
description = "Room 01"

[[sites.buildings.floors.rooms.rows]]
name = "row-01"
description = "Row 01"

[[sites.buildings.floors.rooms.rows.racks]]
name = "rack-01"
description = "Rack 01"
`

func init() {
	RootCmd.AddCommand(newCmd)
}
