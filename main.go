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

	"github.com/robphoenix/go-aci/aci"
)

// NodesJSON ...
type NodesJSON struct {
	FabricNodeIdentPol FabricNodeIdentPol `json:"fabricNodeIdentPol"`
}

// NodeJSON ...
type NodeJSON struct {
	FabricNodeIdentP FabricNodeIdentP `json:"fabricNodeIdentP"`
}

// FabricNodeIdentPol ...
type FabricNodeIdentPol struct {
	Attributes Attributes `json:"attributes"`
	Children   []NodeJSON `json:"children"`
}

// Attributes ...
type Attributes struct {
	Name   string `json:"name,omitempty"`
	NodeID string `json:"nodeId,omitempty"`
	Role   string `json:"role,omitempty"`
	Serial string `json:"serial,omitempty"`
	Status string `json:"status,omitempty"`
}

// FabricNodeIdentP ...
type FabricNodeIdentP struct {
	Attributes Attributes `json:"attributes"`
}

// GetNodes ...
type GetNodes struct {
	Imdata []struct {
		FabricNode struct {
			Attributes struct {
				AdSt             string `json:"adSt"`
				ChildAction      string `json:"childAction"`
				DelayedHeartbeat string `json:"delayedHeartbeat"`
				Dn               string `json:"dn"`
				FabricSt         string `json:"fabricSt"`
				ID               string `json:"id"`
				LcOwn            string `json:"lcOwn"`
				ModTs            string `json:"modTs"`
				Model            string `json:"model"`
				MonPolDn         string `json:"monPolDn"`
				Name             string `json:"name"`
				NameAlias        string `json:"nameAlias"`
				Role             string `json:"role"`
				Serial           string `json:"serial"`
				Status           string `json:"status"`
				UID              string `json:"uid"`
				Vendor           string `json:"vendor"`
				Version          string `json:"version"`
			} `json:"attributes"`
		} `json:"fabricNode"`
	} `json:"imdata"`
	TotalCount string `json:"totalCount"`
}

func main() {

	apicClient, err := aci.NewClient("https://sandboxapicdc.cisco.com/", "admin", "ciscopsdt")
	if err != nil {
		log.Fatal(err)
	}

	// read in CSV data
	csvFile, err := os.Open(".data/fabric_membership.csv")
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
		fmt.Printf("Name = %+v\n", node.FabricNode.Attributes.Name)
	}
}
