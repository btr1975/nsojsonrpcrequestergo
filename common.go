package nsojsonrpcrequestergo

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

// nsoJsonRpcHTTPConnection holds the connection data
type nsoJsonRpcHTTPConnection struct {
	protocol, ip, username, password string
	port                             int
	sslVerify                        bool
	headers                          nsoRequestHeaders
}

// nsoRequestHeaders holds the common request headers
type nsoRequestHeaders struct {
	ContentType string `json:"Content-Type"`
	Accept      string `json:"Accept"`
}

// Constructor to create a new newNsoJsonRpcHTTPConnection struct
//   :values protocol: http, https
//   :values ip: a IPv4 address, or a CNAME
//   :values port: 1 to 65535
//   :values username: A username
//   :values password: A password
//   :values sslVerify: true to verify SSL, false not to
func newNsoJsonRpcHTTPConnection(protocol string, ip string, port int, username string, password string, sslVerify bool) (*nsoJsonRpcHTTPConnection, error) {

	// Check if protocol is http, or https
	if protocol == "http" || protocol == "https" {
	} else {
		return &nsoJsonRpcHTTPConnection{}, errors.New("only http, and https is supported")
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
			return &nsoJsonRpcHTTPConnection{}, err

		}

	}

	if (port < 1) || (port > 65535) {
		return &nsoJsonRpcHTTPConnection{}, errors.New("valid port range between 1 and 65535")
	}

	headers := nsoRequestHeaders{
		ContentType: "application/json",
		Accept:      "application/json",
	}

	return &nsoJsonRpcHTTPConnection{protocol: protocol, ip: ip, port: port, username: username, password: password, sslVerify: sslVerify, headers: headers}, nil

}

// Method to get the NSO JsonRPC URL
func (c *nsoJsonRpcHTTPConnection) NsoUrl() string {
	return fmt.Sprintf("%s://%s:%d/jsonrpc", c.protocol, c.ip, c.port)
}

// Method to get the NSO Headers
func (c *nsoJsonRpcHTTPConnection) NsoHeaders() *nsoRequestHeaders {
	return &c.headers

}

/*
END OF NSO Server connection
*/

/*
START OF NSO JSON-Rpc Requester
*/

type nsoJsonConnection struct {
	request *req.Req
	id      int
	th      float64
	nsocon  nsoJsonRpcHTTPConnection
}

// Constructor to create a new newNsoJsonConnection struct
//   :values protocol: http, https
//   :values ip: a IPv4 address, or a CNAME
//   :values port: 1 to 65535
//   :values username: A username
//   :values password: A password
//   :values sslVerify: true to verify SSL, false not to
func newNsoJsonConnection(protocol string, ip string, port int, username string, password string, sslVerify bool) (*nsoJsonConnection, error) {
	rand.Seed(int64(time.Now().Second()))
	newId := rand.Intn(65000-1+1) + 1

	c, err := newNsoJsonRpcHTTPConnection(protocol, ip, port, username, password, sslVerify)

	if err != nil {
		return &nsoJsonConnection{}, err
	}

	return &nsoJsonConnection{id: newId, nsocon: *c}, nil

}

// Method to convert the NsoJsonRequest to a bytes.Buffer for transport
//   :values param: A req.Param
func (nsoJson *nsoJsonConnection) getJsonRequest(param req.Param) *bytes.Buffer {

	jsonData, _ := json.Marshal(param)

	return bytes.NewBuffer(jsonData)

}

// Method to send a POST request
//   :values param: A req.Param
func (nsoJson *nsoJsonConnection) sendPost(param req.Param) (*req.Resp, error) {
	if nsoJson.nsocon.sslVerify == true {
		nsoJson.request.EnableInsecureTLS(false)

	} else {
		nsoJson.request.EnableInsecureTLS(true)

	}

	response, err := nsoJson.request.Post(nsoJson.nsocon.NsoUrl(), req.BodyJSON(nsoJson.getJsonRequest(param)), req.HeaderFromStruct(nsoJson.nsocon.NsoHeaders()))

	if err != nil {
		return response, err
	}

	return response, nil

}

// Method to send a GET request
//   :values param: A req.Param
func (nsoJson *nsoJsonConnection) sendGet(param req.Param) (*req.Resp, error) {
	if nsoJson.nsocon.sslVerify == true {
		nsoJson.request.EnableInsecureTLS(false)

	} else {
		nsoJson.request.EnableInsecureTLS(true)

	}

	response, err := nsoJson.request.Get(nsoJson.nsocon.NsoUrl(), req.BodyJSON(nsoJson.getJsonRequest(param)), req.HeaderFromStruct(nsoJson.nsocon.NsoHeaders()))

	if err != nil {
		return response, err
	}

	return response, nil

}

// Method to login to the NSO Server
func (nsoJson *nsoJsonConnection) NsoLogin() error {
	param := req.Param{
		"jsonrpc": "2.0",
		"id":      nsoJson.id,
		"method":  "login",
		"params":  map[string]string{"user": nsoJson.nsocon.username, "passwd": nsoJson.nsocon.password},
	}

	request := req.New()
	nsoJson.request = request

	_, err := nsoJson.sendPost(param)

	if err != nil {
		return err
	}

	return nil

}

// Method to logout to the NSO Server
func (nsoJson *nsoJsonConnection) NsoLogout() error {
	param := req.Param{
		"jsonrpc": "2.0",
		"id":      nsoJson.id,
		"method":  "logout",
	}

	response, _ := nsoJson.sendPost(param)

	err := response.Response().Body.Close()

	if err != nil {
		return err
	}

	return nil
}

// Method to start a new NSO Transaction
//   :values mode: read, or read_write
//   :values confMode: private, shared, or exclusive
//   :values tag: "" or a value
//   :values onPendingChanges: reuse, reject, or discard
func (nsoJson *nsoJsonConnection) NewTransaction(mode, confMode, tag, onPendingChanges string) error {
	param := req.Param{
		"jsonrpc": "2.0",
		"id":      nsoJson.id,
		"method":  "new_trans",
		"params": map[string]string{
			"db":                 "running",
			"mode":               mode,
			"conf_mode":          confMode,
			"tag":                tag,
			"on_pending_changes": onPendingChanges,
		},
	}

	response, err := nsoJson.sendPost(param)

	if err != nil {
		return err
	}

	nsoResponse := NewNsoJsonResponse()
	nsoJson.th = nsoResponse.GetTransactionHandle(response)

	return nil
}

// Method to get all NSO transactions
func (nsoJson *nsoJsonConnection) GetTransaction() (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id":      nsoJson.id,
		"method":  "get_trans",
	}

	response, err := nsoJson.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get NSO system settings
//   :values operation: capabilities, customizations , models, user, version, or all
func (nsoJson *nsoJsonConnection) GetSystemSetting(operation string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id":      nsoJson.id,
		"method":  "get_system_setting",
		"params": map[string]string{
			"operation": operation,
		},
	}

	response, err := nsoJson.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to abort a request-id
//   :values requestID: An id
func (nsoJson *nsoJsonConnection) Abort(requestID int) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id":      nsoJson.id,
		"method":  "abort",
		"params": map[string]int{
			"id": requestID,
		},
	}

	response, err := nsoJson.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to evaluate a xpath expression
//   :values xpathExpression: An xpath expression
func (nsoJson *nsoJsonConnection) EvalXPATH(xpathExpression string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id":      nsoJson.id,
		"method":  "eval_xpath",
		"params": map[string]interface{}{
			"th":         nsoJson.th,
			"xpath_expr": xpathExpression,
		},
	}

	response, err := nsoJson.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

/*
END OF NSO JSON-Rpc Requester
*/
