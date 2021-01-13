package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req"
	"math/rand"
	"net"
	"time"
)

/*
START OF NSO Server connection
*/

// NsoJsonRpcHTTPConnection holds the connection data
type NsoJsonRpcHTTPConnection struct {
	protocol, ip, username, password string
	port int
	sslVerify bool
	headers nsoRequestHeaders
}

// nsoRequestHeaders holds the common request headers
type nsoRequestHeaders struct {
	ContentType string `json:"Content-Type"`
	Accept string `json:"Accept"`
}

// Contructor to create a new NsoJsonRpcHTTPConnection struct
func NewNsoJsonRpcHTTPConnection(protocol string, ip string, port int, username string, password string, sslVerify bool) (*NsoJsonRpcHTTPConnection, error)  {

	// Check if protocol is http, or https
	if protocol == "http" || protocol == "https" {
	} else {
		return &NsoJsonRpcHTTPConnection{}, errors.New("Only http, and https is supported!!")
	}

	// Check if ip given is a Ucast IP
	_, err := IpV4UnicastAddress(ip)

	// If err comes back try a DNS lookup
	if err != nil {
		foundlookup, _ := net.LookupIP(ip)

		// If one or more IPs looked up, grab the first found IPv4
		if len(foundlookup) >= 1 {
			for _, v := range foundlookup {
				foundip, _ := IpV4Address(v.String())
				ip = foundip
				break
			}

			// If no lookup found send error
		} else {
			return &NsoJsonRpcHTTPConnection{}, err

		}

	}

	if (port < 1) || (port > 65535) {
		return &NsoJsonRpcHTTPConnection{}, errors.New("valid port range between 1 and 65535")
	}

	headers := nsoRequestHeaders{
		ContentType: "application/json",
		Accept:      "application/json",
	}

	return &NsoJsonRpcHTTPConnection{protocol: protocol, ip: ip, port: port, username: username, password: password, sslVerify: sslVerify, headers: headers}, nil

}

// Method to get the NSO JsonRPC URL
func (c *NsoJsonRpcHTTPConnection) NsoUrl() string {
	return fmt.Sprintf("%s://%s:%d/jsonrpc", c.protocol, c.ip, c.port)
}

// Method to get the NSO Headers
func (c *NsoJsonRpcHTTPConnection) NsoHeaders() *nsoRequestHeaders {
	return &c.headers

}

/*
END OF NSO Server connection
*/

/*
START OF NSO JSON-Rpc Requeser
*/

type NsoJsonConnection struct {
	request *req.Req
	id int
	th int
	nsocon NsoJsonRpcHTTPConnection
}

func NewNsoJsonConnection(c *NsoJsonRpcHTTPConnection) (*NsoJsonConnection, error) {
	rand.Seed(int64(time.Now().Second()))
	newId := rand.Intn(65000 - 1 + 1) + 1

	return &NsoJsonConnection{id: newId, nsocon: *c}, nil

}

// Method to convert the NsoJsonRequest to a bytes.Buffer for transport
func (nsoJson *NsoJsonConnection) getJsonRequest(param req.Param) *bytes.Buffer {

	jsonData, _ := json.Marshal(param)

	return bytes.NewBuffer(jsonData)

}

// Method to send a POST request
func (nsoJson *NsoJsonConnection) sendPost(param req.Param) (*req.Resp, error) {
	response, err := nsoJson.request.Post(nsoJson.nsocon.NsoUrl(), req.BodyJSON(nsoJson.getJsonRequest(param)), req.HeaderFromStruct(nsoJson.nsocon.NsoHeaders()))

	if err != nil {
		return response, err
	}

	return response, nil

}

// Method to send a GET request
func (nsoJson *NsoJsonConnection) sendGet(param req.Param) (*req.Resp, error) {
	response, err := nsoJson.request.Get(nsoJson.nsocon.NsoUrl(), req.BodyJSON(nsoJson.getJsonRequest(param)), req.HeaderFromStruct(nsoJson.nsocon.NsoHeaders()))

	if err != nil {
		return response, err
	}

	return response, nil

}

// Method to login to the NSO Server
func (nsoJson *NsoJsonConnection) NsoLogin(username, password string) *req.Resp {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": nsoJson.id,
		"method": "login",
		"params": map[string]string{"user": username, "passwd": password},
	}

	request := req.New()
	nsoJson.request = request

	response, _ := nsoJson.sendPost(param)

	return response

}

// Method to logout to the NSO Server
func (nsoJson *NsoJsonConnection) NsoLogout() error {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": nsoJson.id,
		"method": "logout",
	}

	response, _ := nsoJson.sendPost(param)

	err := response.Response().Body.Close()

	if err != nil {
		return err
	}

	return nil
}

func (nsoJson *NsoJsonConnection) NewTransaction(mode, confMode, tag, onPendingChanges string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": nsoJson.id,
		"method": "new_trans",
		"params": map[string]string{
			"db": "running",
			"mode": mode,
			"conf_mode": confMode,
			"tag": tag,
			"on_pending_changes": onPendingChanges,
		},
	}

	response, err := nsoJson.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

/*
END OF NSO JSON-Rpc Requeser
*/

// NsoJsonResponse holds a basic NSO JSON RPC Response
// The tags help to convert fields to lowercase
type NsoJsonResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result map[string]interface{} `json:"result"`
	ID int `json:"id"`
	Error map[string]interface{} `json:"error"`
}


