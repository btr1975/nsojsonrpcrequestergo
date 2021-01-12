package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nsojsonrpcrequestergo/common"
)


func main()  {

	nsoConnection, err := common.NewNsoConnection("http", "10.0.0.146", 8080, "admin", "admin", false)

	if err != nil {
		fmt.Println(nsoConnection, err)
	}

	request, _ := common.NewNsoJsonRequest()

	response := request.NsoLogin(nsoConnection)

	defer response.Body.Close()

	fmt.Println(response.Cookies())



	//f, _ := io.Copy(os.Stdout, response.Body)

	newTemp := common.NsoJsonResponse{}



	z, _ := ioutil.ReadAll(response.Body)

	_ = json.Unmarshal(z, &newTemp)

	fmt.Println(newTemp)

	fmt.Println(string(z))


	//fmt.Println(string(y))







}

