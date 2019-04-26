package main

/*
Reading data form XML string listed below

<Data>
	<Node1>100</Node1>
	<Node2>200</Node2>
</Data>

Found nice app which create Golang structures from and XML file
use followed intallation command:

go get github.com/gnewton/chidley
*/

import (
	"fmt"
	xml "encoding/xml"
)

// Data represents and XML reflection
type Data struct {
	XMLName xml.Name `xml:"Data"`
	Node1 int `xml:"Node1,omitempty"`
	Node2 *int `xml:"Node2,omitempty"`
}																	

func main() {
	var d Data
	
	err := xml.Unmarshal([]byte("<Data><Node1>100</Node1><Node2>250</Node2></Data>"), &d)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if d.Node2 == nil {
		fmt.Println("No Node2 existed in XML")
		return
	}
	fmt.Printf("Node2 value is %d\n", *d.Node2)
}