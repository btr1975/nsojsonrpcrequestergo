package main

import (
	"fmt"
	"nsojsonrpcrequestergo/common"
)


func main()  {


	nsoHTTPConnection, err := common.NewNsoJsonRpcHTTPConnection("http", "10.0.0.146", 8080, "admin", "admin", false)

	if err != nil {
		fmt.Println(nsoHTTPConnection, err)
	}

	nsoConnection, err := common.NewNsoJsonConnection(nsoHTTPConnection)

	if err != nil {
		fmt.Println(nsoConnection, err)
	}

	err = nsoConnection.NsoLogin("admin", "admin")

	if err != nil {
		fmt.Println(err)
	}

	err = nsoConnection.NewTransaction("read", "private", "", "reuse")

	if err != nil {
		fmt.Println(err)
	}

    config, err := common.NewNsoJsonRpcConfig(nsoConnection)

	if err != nil {
		fmt.Println(err)
	}

    selections := []string{"device-name", "device-type"}
    var sort []string

    nsoQuery, err := common.NewQueryObject("/services/etradeing_spine_and_leaf_devices", "", selections,0, 0, sort,  "", true, "", "string")

	if err != nil {
		fmt.Println(err)
	}

	err = config.StartQuery(nsoQuery)

	if err != nil {
		fmt.Println(err)
	}

	nsoQueryData, _ := config.RunQuery(nsoQuery)

	data := common.NewNsoJsonResponse()

	results, err := data.GetQueryResults(nsoQueryData)

	if err != nil {
		panic("no")
	}

	fmt.Println(results)

	err = config.StopQuery(nsoQuery)

	if err != nil {
		fmt.Println(err)
	}


	err = nsoConnection.NsoLogout()

	if err != nil {
		fmt.Println(err)
	}



}
