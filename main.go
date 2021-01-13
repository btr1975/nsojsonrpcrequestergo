package main

import (
	"fmt"
	"nsojsonrpcrequestergo/common"
)


func main()  {

	nsoConnection, err := common.NewNsoJsonRpcHTTPConnection("http", "10.0.0.146", 8080, "admin", "admin", false)

	if err != nil {
		fmt.Println(nsoConnection, err)
	}

	// Using req lib

	// req.Debug = true

	thing, _ := common.NewNsoJsonConnection(nsoConnection)

	_ = thing.NsoLogin("admin", "admin")

	thing2, _ := thing.NewTransaction("read", "private", "", "reuse")

	fmt.Println(thing2)

	thing3, _ := thing.GetSystemSetting("all")

	fmt.Println(thing3)


	err = thing.NsoLogout()

	fmt.Println(err)










}

