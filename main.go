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

// NodesJSON ...
type NodesJSON struct {
	FabricNodeIdentPol `json:"fabricNodeIdentPol"`
}

// NodeJSON ...
type NodeJSON struct {
	FabricNodeIdentP `json:"fabricNodeIdentP"`
}

// FabricNodeIdentPol ...
type FabricNodeIdentPol struct {
	Attributes `json:"attributes"`
	Children   []NodeJSON `json:"children"`
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

// FabricNodeIdentP ...
type FabricNodeIdentP struct {
	Attributes `json:"attributes"`
}

// GetNodes ...
type GetNodes struct {
	Imdata []struct {
		FabricNode `json:"fabricNode"`
	} `json:"imdata"`
	TotalCount string `json:"totalCount"`
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
	// fmt.Printf("%v\n", nodes)

	// instantiate a nodes struct
	fnipol := FabricNodeIdentPol{
		Attributes: Attributes{
			Status: "created,modified",
		},
	}
	// add individual nodes to the nodes struct
	for _, n := range nodes {
		fnip := NodeJSON{
			FabricNodeIdentP: FabricNodeIdentP{
				Attributes: Attributes{
					Name:   n["Name"],
					NodeID: n["Node ID"],
					Role:   n["Role"],
					Serial: n["Serial"],
					// create the node
					Status: "created,modified",
					// delete the node
					// Status: "deleted",
				},
			}}
		fnipol.Children = append(fnipol.Children, fnip)
		// fmt.Printf("%q\n", fnip)
	}
	nj := NodesJSON{fnipol}
	// marshal the struct into JSON
	b, err := json.Marshal(nj)
	if err != nil {
		fmt.Println("error:", err)
	}

	// login
	err = apicClient.Login()
	if err != nil {
		log.Fatal(err)
	}

	// nodes endpoint
	NodesURL := "https://sandboxapicdc.cisco.com/api/node/mo/uni/controller/nodeidentpol.json"
	req, err := http.NewRequest("POST", NodesURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", apicClient.Cookie)
	client := apicClient.Client
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	nodesBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(nodesBody))

	// client = &http.Client{Transport: T}
	URL := "https://sandboxapicdc.cisco.com/api/node/class/fabricNode.json"
	req, err = http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Cookie", apicClient.Cookie)
	resp, err = client.Do(req)
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
	for _, node := range n.Imdata {
		fmt.Printf("Name = %+v\n", node.Name)
	}
}
