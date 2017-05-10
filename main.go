package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
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
					// Status: "created,modified",
					Status: "deleted",
				},
			}}
		fnipol.Children = append(fnipol.Children, fnip)
		// fmt.Printf("%q\n", fnip)
	}
	nj := NodesJSON{fnipol}
	b, err := json.Marshal(nj)
	if err != nil {
		fmt.Println("error:", err)
	}

	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       true,
			MinVersion:               tls.VersionTLS11,
			MaxVersion:               tls.VersionTLS11,
		},
	}

	loginURL := "https://sandboxapicdc.cisco.com/api/aaaLogin.json"
	loginString := fmt.Sprintf(`{"aaaUser": {"attributes": {"name": "admin", "pwd": "ciscopsdt"}}}`)
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer([]byte(loginString)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Transport: t}
	resp, err := client.Do(req)
	cookies := resp.Cookies()
	apicCookie := cookies[0]
	token := apicCookie.Value
	tokenName := apicCookie.Name
	fmt.Println("Token: ", token)
	fmt.Println("Token Name: ", tokenName)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	loginBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(loginBody))

	NodesURL := "https://sandboxapicdc.cisco.com/api/node/mo/uni/controller/nodeidentpol.json"
	req, err = http.NewRequest("POST", NodesURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", apicCookie.String())
	client = &http.Client{Transport: t}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	nodesBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(nodesBody))
}
