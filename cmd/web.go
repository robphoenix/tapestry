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
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/apic", apicHandler)
	http.HandleFunc("/fabric-membership", fabricMembershipHandler)
	http.HandleFunc("/geolocation", geolocationHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	base := filepath.Join("templates", "base.html.tmpl")
	content := filepath.Join("templates", "index.html.tmpl")
	var tmpl = template.Must(template.ParseFiles(base, content))
	if err := tmpl.Execute(w, ""); err != nil {
		log.Println(err)
	}

}
func apicHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	switch r.Method {
	case "GET":
		base := filepath.Join("templates", "base.html.tmpl")
		content := filepath.Join("templates", "apic.html.tmpl")
		var tmpl = template.Must(template.ParseFiles(base, content))
		if err := tmpl.Execute(w, cfg.APIC); err != nil {
			log.Println(err)
		}
	case "POST":
		r.ParseForm()
		vals := r.PostForm
		url := vals["url"]
		user := vals["username"]
		pass := vals["password"]
		cfg.APIC.URL = url[0]
		cfg.APIC.Username = user[0]
		cfg.APIC.Password = pass[0]

		b, err := toml.Marshal(cfg)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(viper.ConfigFileUsed(), b, os.FileMode(0644))
		if err != nil {
			log.Fatal(err)
		}
		base := filepath.Join("templates", "base.html.tmpl")
		content := filepath.Join("templates", "apic.html.tmpl")
		var tmpl = template.Must(template.ParseFiles(base, content))
		if err := tmpl.Execute(w, cfg.APIC); err != nil {
			log.Println(err)
		}
	}
}

func fabricMembershipHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	base := filepath.Join("templates", "base.html.tmpl")
	content := filepath.Join("templates", "fabric_membership.html.tmpl")
	var tmpl = template.Must(template.ParseFiles(base, content))
	if err := tmpl.Execute(w, cfg.Nodes); err != nil {
		log.Println(err)
	}
}

func geolocationHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	base := filepath.Join("templates", "base.html.tmpl")
	content := filepath.Join("templates", "geolocation.html.tmpl")
	var tmpl = template.Must(template.ParseFiles(base, content))
	if err := tmpl.Execute(w, cfg.Sites); err != nil {
		log.Println(err)
	}
}

func init() {
	RootCmd.AddCommand(webCmd)
}
