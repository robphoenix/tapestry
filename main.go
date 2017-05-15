package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/robphoenix/go-aci/aci"
)

// GetNodes ...
type GetNodes struct {
	Imdata []struct {
		FabricNode `json:"fabricNode"`
	} `json:"imdata"`
	TotalCount string `json:"totalCount"`
}

// Attributes ...
type Attributes struct {
	Name             string `json:"name,omitempty"`
	NodeID           string `json:"nodeId,omitempty"`
	Role             string `json:"role,omitempty"`
	Serial           string `json:"serial,omitempty"`
	Status           string `json:"status,omitempty"`
	AdSt             string `json:"adSt,omitempty"`
	ChildAction      string `json:"childAction,omitempty"`
	DelayedHeartbeat string `json:"delayedHeartbeat,omitempty"`
	Dn               string `json:"dn,omitempty"`
	FabricSt         string `json:"fabricSt,omitempty"`
	ID               string `json:"id,omitempty"`
	LcOwn            string `json:"lcOwn,omitempty"`
	ModTs            string `json:"modTs,omitempty"`
	Model            string `json:"model,omitempty"`
	MonPolDn         string `json:"monPolDn,omitempty"`
	NameAlias        string `json:"nameAlias,omitempty"`
	UID              string `json:"uid,omitempty"`
	Vendor           string `json:"vendor,omitempty"`
	Version          string `json:"version,omitempty"`
}

// FabricNode ...
type FabricNode struct {
	Attributes `json:"attributes"`
}

func main() {

	// cmd.Execute()
	config, err := toml.LoadFile("tapestry.toml")
	if err != nil {
		log.Fatal(err)
	}
	user := config.Get("apic.username").(string)
	password := config.Get("apic.password").(string)
	url := config.Get("apic.url").(string)
	dataSrc := config.Get("data.src").(string)
	fabricNodeSrc := config.Get("fabricNodes.src").(string)

	apicClient, err := aci.NewClient(url, user, password)
	if err != nil {
		log.Fatal(err)
	}

	// read in CSV data
	fabricNodesDataFile := filepath.Join(dataSrc, fabricNodeSrc)
	csvFile, err := os.Open(fabricNodesDataFile)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	// headers
	headers, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	// records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// build a slice of nodes, a node is a map of headers to records
	// struct??
	var nodes []map[string]string
	for _, r := range records {
		node := make(map[string]string, len(headers))
		for i := 0; i < len(headers); i++ {
			node[headers[i]] = r[i]
		}
		nodes = append(nodes, node)
	}

	var ns []aci.Node
	for _, node := range nodes {
		n := aci.Node{
			Name:   node["Name"],
			ID:     node["Node ID"],
			Role:   node["Role"],
			Serial: node["Serial"],
		}
		ns = append(ns, n)
	}

	// login
	err = apicClient.Login()
	if err != nil {
		log.Fatal(err)
	}

	// add nodes
	err = apicClient.AddNodes(ns)
	if err != nil {
		log.Fatal(err)
	}

	// list nodes
	URL := "https://sandboxapicdc.cisco.com/api/node/class/fabricNode.json"
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Cookie", apicClient.Cookie)
	resp, err := apicClient.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	nodesList, _ := ioutil.ReadAll(resp.Body)
	var n GetNodes
	err = json.NewDecoder(bytes.NewReader(nodesList)).Decode(&n)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("response Body:", string(nodesList))
	// fmt.Printf("%#v\n", n)
	for _, node := range n.Imdata {
		fmt.Printf("Name = %+v\n", node.Name)
	}
}
