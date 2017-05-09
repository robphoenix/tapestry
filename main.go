package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type NodesJSON struct {
	FabricNodeIdentPol FabricNodeIdentPol `json:"fabricNodeIdentPol"`
}

type NodeJSON struct {
	FabricNodeIdentP FabricNodeIdentP `json:"fabricNodeIdentP"`
}

type FabricNodeIdentPol struct {
	Attributes Attributes `json:"attributes"`
	Children   []NodeJSON `json:"children"`
}

type Attributes struct {
	Name   string `json:"name,omitempty"`
	NodeID string `json:"nodeId,omitempty"`
	Role   string `json:"role,omitempty"`
	Serial string `json:"serial,omitempty"`
	Status string `json:"status,omitempty"`
}

type FabricNodeIdentP struct {
	Attributes Attributes `json:"attributes"`
}

func main() {
	csvFile, err := os.Open(".data/fabric_membership.csv")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	headers, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var nodes []map[string]string

	for _, r := range records {
		node := make(map[string]string, len(headers))
		for i := 0; i < len(headers); i++ {
			node[headers[i]] = r[i]
		}
		nodes = append(nodes, node)
	}
	// fmt.Printf("%v\n", nodes)

	fnipol := FabricNodeIdentPol{
		Attributes: Attributes{
			Status: "created,modified",
		},
	}
	for _, n := range nodes {
		fnip := NodeJSON{
			FabricNodeIdentP: FabricNodeIdentP{
				Attributes: Attributes{
					Name:   n["Name"],
					NodeID: n["Node ID"],
					Role:   n["Role"],
					Serial: n["Serial"],
					Status: "created,modified",
				},
			}}
		fnipol.Children = append(fnipol.Children, fnip)
		fmt.Printf("%q\n", fnip)
	}
	// fmt.Printf("%q\n", fnipol)
	// for _, v := range fnipol.Children {
	//         fmt.Println(v.Attributes)
	// }
	nj := NodesJSON{fnipol}
	b, err := json.Marshal(nj)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}
