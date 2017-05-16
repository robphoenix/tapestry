package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/robphoenix/go-aci/aci"
)

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

	// // delete nodes
	// err = apicClient.DeleteNodes(ns)
	// if err != nil {
	//         log.Fatal(err)
	// }

	// list nodes
	n, err := apicClient.ListNodes()
	// fmt.Printf("%#v\n", n)
	for _, node := range n {
		fmt.Printf("node is an %T type\n", node)
		fmt.Printf("Name = %+v\n", node.Name)
		fmt.Printf("ID = %+v\n", node.ID)
		fmt.Printf("Serial = %+v\n", node.Serial)
		fmt.Printf("fabric status = %+v\n", node.FabricStatus)
		fmt.Printf("node.Status = %+v\n", node.Status)
	}
}
