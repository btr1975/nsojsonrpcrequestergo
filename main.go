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

    selections := []string{"device-name", "device-type"}
    var sort []string

    queryData, _ := common.NewQueryObject("/services/etradeing_spine_and_leaf_devices", "", selections,0, 0, sort,  "", true, "", "string")

	thing3, _ := config.StartQuery(queryData)
	thing4, _ := config.RunQuery(queryData)

	data := common.NewNsoJsonResponse()

	newData, _ := data.ResponseToStruct(thing4)

	resutlData := newData.Result

	for k, v := range resutlData {
		if k == "results" {
			fmt.Println(k)
			fmt.Println(v)
		}
	}


	fmt.Println(resutlData)

	thing5, _ := config.StopQuery(queryData)



	fmt.Println(thing3)
	fmt.Println(thing4)
	fmt.Println(thing5)

	err = thing.NsoLogout()

	fmt.Println(err)










}

