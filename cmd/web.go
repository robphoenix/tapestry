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
	"io/ioutil"
	"log"
	"net/http"
	"os"

	toml "github.com/pelletier/go-toml"
	"github.com/robphoenix/tapestry/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// webCmd represents the editor command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Edit your Tapestry project configuration.",
	Long:  `A web UI to visually edit your Tapestry project configuration.`,
	Run:   runWeb,
}

func runWeb(cmd *cobra.Command, args []string) {
	http.HandleFunc("/", webHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	switch r.Method {
	case "GET":
		if err := webTemplate.Execute(w, cfg); err != nil {
			log.Println(err)
		}
	case "POST":
		r.ParseForm()
		vals := r.PostForm
		if len(vals["apicSubmit"]) > 0 {
			cfg.APIC.URL = vals["url"][0]
			cfg.APIC.Username = vals["username"][0]
			cfg.APIC.Password = vals["password"][0]
		}
		if len(vals["fabricMembershipSubmit"]) > 0 {
			var ns []config.Node
			for i := 0; i < len(vals["name"]); i++ {
				ns = append(ns, config.Node{
					ID:     vals["id"][i],
					Name:   vals["name"][i],
					Pod:    vals["pod"][i],
					Serial: vals["serial"][i],
					Role:   vals["role"][i],
				})
			}
			cfg.Nodes = ns
		}

		b, err := toml.Marshal(cfg)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(viper.ConfigFileUsed(), b, os.FileMode(0644))
		if err != nil {
			log.Fatal(err)
		}
		if err := webTemplate.Execute(w, cfg); err != nil {
			log.Println(err)
		}
	}

}

func init() {
	RootCmd.AddCommand(webCmd)
}
