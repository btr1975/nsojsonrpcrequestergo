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





func main()  {

	thing, err := common.NewCommon("http", "mainpc.tsnetsolutions.local", 65535, "admin", "admin", false)

	if err != nil {
		fmt.Println(thing, err)
	}

	fmt.Println(thing, err)


	jsond, _ := common.NewNosJsonRequest(20, "login", map[string]string{"user": "admin", "passwd": "admin"})

	a, _ := json.Marshal(jsond)

	client  := &http.Client{Timeout: 60}
	req, _ := http.NewRequest("POST", "http://10.0.0.146:8080/jsonrpc", bytes.NewBuffer(a))
	req.Header.Add("Content-Type","application/json")
	req.Header.Add("Accept","application/json")

	response, _ := client.Do(req)

	fmt.Println(response.Cookies())

	f, _ := io.Copy(os.Stdout, response.Body)

	fmt.Println(string(f))





}

