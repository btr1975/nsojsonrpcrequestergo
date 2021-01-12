package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nsojsonrpcrequestergo/common"
	"os"
)


type JsonRequest struct {
	Jsonrpc, id, method string
	params map[string]string
}


func main()  {

	thing, err := common.NewCommon("http", "mainpc.tsnetsolutions.local", 65535, "admin", "admin", false)

	if err != nil {
		fmt.Println(thing, err)
	}

	fmt.Println(thing, err)


	jsond := &JsonRequest{
		Jsonrpc: "2.0",
		id:      "1000",
		method:  "login",
		params: map[string]string{"user": "admin", "passwd": "admin"},

	}

	a, _ := json.Marshal(jsond)

	fmt.Println(string(a))

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(jsond)
	fmt.Println(jsond)


	client  := &http.Client{}
	req, _ := http.NewRequest("POST", "http://10.0.0.146:8080/jsonrpc", b)
	req.Header.Add("Content-Type","application/json")
	req.Header.Add("Accept","application/json")

	response, _ := client.Do(req)

	fmt.Println(response)

	io.Copy(os.Stdout, response.Body)





}

