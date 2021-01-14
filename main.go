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

    config, _ := common.NewNsoJsonRpcConfig(thing)


	thing3, _ := config.Query("/services/etradeing_spine_and_leaf_devices[device-name='UNIT-TEST-NX-LEA10']/device-type", "keypath-value")

	fmt.Println(thing3)

	err = thing.NsoLogout()

	fmt.Println(err)










}

