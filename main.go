package main

import (
	"fmt"
	"nsojsonrpcrequestergo/common"
)



func main()  {

	thing, err := common.NewCommon("http", "mainpc.tsnetsolutions.local", 65535, "admin", "admin", false)

	if err != nil {
		fmt.Println(thing, err)
	}

	fmt.Println(thing, err)

}

