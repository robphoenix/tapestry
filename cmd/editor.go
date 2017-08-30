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
	"net/http"

	"github.com/spf13/cobra"
)

// editorCmd represents the editor command
var editorCmd = &cobra.Command{
	Use:   "editor",
	Short: "Edit your Tapestry project configuration.",
	Long:  `A web UI to visually edit your Tapestry project configuration.`,
	Run:   runEditor,
}

func runEditor(cmd *cobra.Command, args []string) {
	fmt.Println("editor called")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/apic", apicHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tapestry.")
}

func apicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%+v", cfg.APIC)
}

func init() {
	RootCmd.AddCommand(editorCmd)
}
