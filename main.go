package main

import (
	"fmt"
	"nsojsonrpcrequestergo/common"
	"strings"
)


func fixResult(result string) []string {
	remLeftDblBracket := strings.Replace(result, "[[","", -1)
	remRightDblBracket := strings.Replace(remLeftDblBracket, "]]","", -1)
	csvFmt := strings.Replace(remRightDblBracket, "] [",",", -1)
	return strings.Split(csvFmt, ",")

}

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

	results, err := data.GetQueryResults(thing4)

	if err != nil {
		panic("no")
	}

	fmt.Println(results)



	thing5, _ := config.StopQuery(queryData)



	fmt.Println(thing3)
	fmt.Println(thing4)
	fmt.Println(thing5)

	err = thing.NsoLogout()

	fmt.Println(err)

}
