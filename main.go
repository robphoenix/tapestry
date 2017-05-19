package main

import (
	"github.com/robphoenix/tapestry/cmd"
)

func main() {

	cmd.Execute()

	// // fetch configuration data
	// conf, err := newConfig()
	// if err != nil {
	//         log.Fatal(err)
	// }
	//
	// // create new APIC client
	// apicClient, err := aci.NewClient(conf.url, conf.user, conf.password)
	// if err != nil {
	//         log.Fatal(err)
	// }
	//
	// // read in data from fabric membership file
	// fabricNodesDataFile := filepath.Join(conf.dataSrc, conf.fabricNodeSrc)
	// nodes, err := newNodes(fabricNodesDataFile)
	// if err != nil {
	//         log.Fatal(err)
	// }
	//
	// // determine desired node state
	// var desiredNodeState []aci.Node
	// for _, node := range nodes {
	//         n := aci.Node{
	//                 Name:   node.Name,
	//                 ID:     node.NodeID,
	//                 Serial: node.Serial,
	//         }
	//         desiredNodeState = append(desiredNodeState, n)
	// }
	//
	// // login
	// err = apicClient.Login()
	// if err != nil {
	//         log.Fatal(err)
	// }
	//
	// // determine actual node state
	// actualNodeState, err := apicClient.ListNodes()
	// if err != nil {
	//         log.Fatal(err)
	// }
	//
	// // determine actions to take
	// na := diffNodeStates(desiredNodeState, actualNodeState)
	//
	// fmt.Printf("%s\n", "Nodes to add:")
	// for _, v := range na.add {
	//         fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	// }
	// fmt.Printf("%s\n", "Nodes to delete:")
	// for _, v := range na.delete {
	//         fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	// }
	//
	// // delete nodes
	// err = apicClient.DeleteNodes(na.delete)
	// if err != nil {
	//         log.Fatal(err)
	// }
	// if err == nil {
	//         fmt.Printf("%s\n", "Nodes deleted:")
	//         for _, v := range na.delete {
	//                 fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	//         }
	// }
	//
	// // add nodes
	// err = apicClient.AddNodes(na.add)
	// if err != nil {
	//         log.Fatal(err)
	// }
	// if err == nil {
	//         fmt.Printf("%s\n", "Nodes added:")
	//         for _, v := range na.add {
	//                 fmt.Printf("%s\t%s\t%s\n", v.Name, v.ID, v.Serial)
	//         }
	// }
}
