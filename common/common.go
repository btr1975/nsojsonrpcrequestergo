package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// NosJsonRequest holds a basic NSO JSON RPC Request
// The tags help to convert fields to lowercase
type NsoJsonRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	ID int `json:"id"`
	Method string `json:"method"`
	Params params `json:"params"`
}

// NsoJsonResponse holds a basic NSO JSON RPC Response
// The tags help to convert fields to lowercase
type NsoJsonResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result map[string]interface{} `json:"result"`
	ID int `json:"id"`
	Error map[string]interface{} `json:"error"`
}

// Struc to hold JSON params
type params map[string]string

// Contructor to create a new NsoJsonRequest struct
func NewNsoJsonRequest() (*NsoJsonRequest, error) {

	rand.Seed(int64(time.Now().Second()))
	thing := rand.Intn(65000 - 1 + 1) + 1

	return &NsoJsonRequest{Jsonrpc: "2.0", ID: thing}, nil

}

func (n *NsoJsonRequest) nsoJsonRequestMethodParmas(method string, params map[string]string) error {
	n.Method = method
	n.Params = params
	return nil
}

// Method to convert the NsoJsonRequest to a bytes.Buffer for transport
func (n *NsoJsonRequest) getJsonRequest() *bytes.Buffer {

	jsonData, _ := json.Marshal(n)

	return bytes.NewBuffer(jsonData)

}


// Method to send a POST login request to NSO
func (n *NsoJsonRequest) NsoLogin(c *NsoConnection) *http.Response  {
	n.Method = "login"
	n.Params = params{"user": c.username, "passwd": c.password}

	client  := &http.Client{}
	req, _ := http.NewRequest("POST", c.NsoUrl(), n.getJsonRequest())
	for k, v := range c.NsoHeaders() {
		req.Header.Add(k,v)

	}

	response, _ := client.Do(req)

	if response.Status != "200 OK" {
		fmt.Println("not 200 OK")
	}

	return response

}



// Method to send a POST request to NSO
func (n *NsoJsonRequest) RequestPost(nsoURL string, nsoHeaders map[string]string) *http.Response  {

	client  := &http.Client{}
	req, _ := http.NewRequest("POST", nsoURL, n.getJsonRequest())
	for k, v := range nsoHeaders {
		req.Header.Add(k,v)

	}

	response, _ := client.Do(req)

	if response.Status != "200 OK" {
		fmt.Println("not 200 OK")
	}

	return response

}



/*
This section starts the NSO Server connection
 */


// NsoConnection holds the connection data
type NsoConnection struct {
	protocol, ip, username, password string
	port int
	sslVerify bool
	headers nsoRequestHeaders
	cookies http.CookieJar
}

// nsoRequestHeaders holds the common request headers
type nsoRequestHeaders map[string]string

// Contructor to create a new NsoConnection struct
func NewNsoConnection(protocol string, ip string, port int, username string, password string, sslVerify bool) (*NsoConnection, error)  {

	// Check if protocol is http, or https
	if protocol == "http" || protocol == "https" {
	} else {
		return &NsoConnection{}, errors.New("Only http, and https is supported!!")
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
			return &NsoConnection{}, err

		}

	}

	if (port < 1) || (port > 65535) {
		return &NsoConnection{}, errors.New("valid port range between 1 and 65535")
	}

	headers := nsoRequestHeaders{"Content-Type": "application/json", "Accept": "application/json"}

	return &NsoConnection{protocol: protocol, ip: ip, port: port, username: username, password: password, sslVerify: sslVerify, headers: headers}, nil

}

// Method to get the NSO JsonRPC URL
func (c *NsoConnection) NsoUrl() string {
	return fmt.Sprintf("%s://%s:%d/jsonrpc", c.protocol, c.ip, c.port)
}

// Method to get the NSO Headers
func (c *NsoConnection) NsoHeaders() map[string]string {
	return c.headers

}

// Method to set the NSO Cookies
func (c *NsoConnection) setCookie(cookies http.CookieJar) error {
	c.cookies = cookies
	return nil
}
