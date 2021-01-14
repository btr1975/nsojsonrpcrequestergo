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

// Constructor to create a new NsoJsonRpcHTTPConnection struct
//   :values protocol: http, https
//   :values ip: a IPv4 address, or a CNAME
//   :values port: 1 to 65535
//   :values username: A username
//   :values password: A password
//   :values sslVerify: true to verify SSL, false not to
func NewNsoJsonRpcHTTPConnection(protocol string, ip string, port int, username string, password string, sslVerify bool) (*NsoJsonRpcHTTPConnection, error)  {

	// Check if protocol is http, or https
	if protocol == "http" || protocol == "https" {
	} else {
		return &NsoJsonRpcHTTPConnection{}, errors.New("only http, and https is supported")
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
START OF NSO JSON-Rpc Requester
*/

type NsoJsonConnection struct {
	request *req.Req
	id int
	th float64
	nsocon NsoJsonRpcHTTPConnection
}

// Constructor to create a new NewNsoJsonConnection struct
//   :values c: A NsoJsonRpcHTTPConnection
func NewNsoJsonConnection(c *NsoJsonRpcHTTPConnection) (*NsoJsonConnection, error) {
	rand.Seed(int64(time.Now().Second()))
	newId := rand.Intn(65000 - 1 + 1) + 1

	return &NsoJsonConnection{id: newId, nsocon: *c}, nil

}

// Method to convert the NsoJsonRequest to a bytes.Buffer for transport
//   :values param: A req.Param
func (nsoJson *NsoJsonConnection) getJsonRequest(param req.Param) *bytes.Buffer {

	jsonData, _ := json.Marshal(param)

	return bytes.NewBuffer(jsonData)

}

// Method to send a POST request
//   :values param: A req.Param
func (nsoJson *NsoJsonConnection) sendPost(param req.Param) (*req.Resp, error) {
	response, err := nsoJson.request.Post(nsoJson.nsocon.NsoUrl(), req.BodyJSON(nsoJson.getJsonRequest(param)), req.HeaderFromStruct(nsoJson.nsocon.NsoHeaders()))

	if err != nil {
		return response, err
	}

	return response, nil

}

// Method to send a GET request
//   :values param: A req.Param
func (nsoJson *NsoJsonConnection) sendGet(param req.Param) (*req.Resp, error) {
	response, err := nsoJson.request.Get(nsoJson.nsocon.NsoUrl(), req.BodyJSON(nsoJson.getJsonRequest(param)), req.HeaderFromStruct(nsoJson.nsocon.NsoHeaders()))

	if err != nil {
		return response, err
	}

	return response, nil

}

// Method to login to the NSO Server
//   :values username: A username
//   :values password: A password
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

// Method to start a new NSO Transaction
//   :values mode: read, or read_write
//   :values confMode: private, shared, or exclusive
//   :values tag: "" or a value
//   :values onPendingChanges: reuse, reject, or discard
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

	nsoResponse := NewNsoJsonResponse()
	nsoJson.th = nsoResponse.GetTransactionHandle(response)

	return response, nil
}

// Method to get all NSO transactions
func (nsoJson *NsoJsonConnection) GetTransaction() (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": nsoJson.id,
		"method": "get_trans",
	}

	response, err := nsoJson.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get NSO system settings
//   :values operation: capabilities, customizations , models, user, version, or all
func (nsoJson *NsoJsonConnection) GetSystemSetting(operation string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": nsoJson.id,
		"method": "get_system_setting",
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
func (nsoJson *NsoJsonConnection) Abort(requestID int) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": nsoJson.id,
		"method": "abort",
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
func (nsoJson *NsoJsonConnection) EvalXPATH(xpathExpression string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": nsoJson.id,
		"method": "eval_xpath",
		"params": map[string]interface{}{
			"th": nsoJson.th,
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

/*
START OF NSO JSON-Rpc Response
*/

// NsoJsonResponse holds a basic NSO JSON RPC Response
// The tags help to convert fields to lowercase
type NsoJsonResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result map[string]interface{} `json:"result"`
	ID int `json:"id"`
	Error map[string]interface{} `json:"error"`
}

func NewNsoJsonResponse() *NsoJsonResponse  {

	return &NsoJsonResponse{}

}

func (r *NsoJsonResponse) GetTransactionHandle(response *req.Resp) float64  {
	var th float64

	_ = response.ToJSON(&r)

	for key, value := range r.Result {
		if key == "th" {
			th = value.(float64)
			break
		}
	}

	return th

}

/*
END OF NSO JSON-Rpc Response
*/

/*
START OF NSO JSON-Rpc Config
*/

type NsoJsonRpcConfig struct {
	nsocon NsoJsonConnection

}

func NewNsoJsonRpcConfig(nsoJson *NsoJsonConnection) (*NsoJsonRpcConfig, error)  {

	return &NsoJsonRpcConfig{nsocon: *nsoJson}, nil
	
}

// Method to show NSO config
//   :values path: A key path
//   :values resultAs: string, or json
//   :values withOper: true for operational data false for not
//   :values maxSize: 0 to disable limit any other number to limit
func (config *NsoJsonRpcConfig) ShowConfig(path, resultAs string, withOper bool, maxSize int) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "show_config",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,
			"result_as": resultAs,
			"with_oper": withOper,
			"max_size": maxSize,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to deref NSO config
//   :values path: A key path
//   :values resultAs: paths, target, or list-target
func (config *NsoJsonRpcConfig) Deref(path, resultAs string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "deref",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,
			"result_as": resultAs,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get leaf reference values
//   :values path: A key path
//   :values skipGrouping: true to skip grouping false to not
//   :values keys: array of keys
func (config *NsoJsonRpcConfig) GetLeafrefValues(path string, skipGrouping bool, keys []string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "get_leafref_values",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,
			"skip_grouping": skipGrouping,
			"keys": keys,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to run an action
//   :values path: A key path
//   :values inputData: A map of data
func (config *NsoJsonRpcConfig) RunAction(path string, inputData map[string]interface{}) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "run_action",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,
			"params": inputData,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get a schema
//   :values path: A key path
func (config *NsoJsonRpcConfig) GetSchema(path string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "get_schema",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get a list of keys
//   :values path: A key path
func (config *NsoJsonRpcConfig) GetListKeys(path string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "get_list_keys",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get a leaf value
//   :values path: A key path
//   :values checkDefault: true to check for default value false to not
func (config *NsoJsonRpcConfig) GetValue(path string, checkDefault bool) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "get_value",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,
			"check_default": checkDefault,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get multiple leaf values
//   :values path: A key path
//   :values leafs: A array of leafs
//   :values checkDefault: true to check for default value false to not
func (config *NsoJsonRpcConfig) GetValues(path string, leafs []string, checkDefault bool) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "get_values",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,
			"check_default": checkDefault,
			"leafs": leafs,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to create a leaf
//   :values path: A key path
func (config *NsoJsonRpcConfig) Create(path string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "create",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}


// Method to check if a leaf exists
//   :values path: A key path
func (config *NsoJsonRpcConfig) Exists(path string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "exists",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get a choice/case
//   :values path: A key path
//   :values choice: A choice from a case
func (config *NsoJsonRpcConfig) GetCase(path, choice string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "get_case",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,
			"choice": choice,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to load data to NSO
//   :values data: The data to be loaded
//   :values path: A key path use "/" at the very least
//   :values dataFormat: json, or xml
//   :values mode: create, merge, or replace
func (config *NsoJsonRpcConfig) Load(data, path, dataFormat, mode string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "load",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"data": data,
			"path": path,
			"format": dataFormat,
			"mode": mode,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to set a value
//   :values path: A key path
//   :values value: What you want to set
//   :values dryRun: true for dryrun false for not
func (config *NsoJsonRpcConfig) SetValue(path string, value interface{}, dryRun bool) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "set_value",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,
			"value": value,
			"dryrun": dryRun,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to validate a commit
//    In the CLI commits are validated automatically, in JsonRPC
//    they are not, but only validated commits can be committed
func (config *NsoJsonRpcConfig) ValidateCommit() (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "validate_commit",
		"params": map[string]interface{}{
			"th": config.nsocon.th,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to commit
//   :values dryRun: true for dryrun false for not
//   :values output: cli, native, or xml
//   :values reverse: true for reverse diff false for forward diff
//                    config only can be used with native
func (config *NsoJsonRpcConfig) Commit(dryRun bool, output string, reverse bool) (*req.Resp, error) {
	var flags []string = nil

	if dryRun == true {
		flags = append(flags, fmt.Sprintf("dry-run=%s", output))
		if output == "native" && reverse == true {
			flags = append(flags, "dry-run-reverse")
		}

	}

	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "commit",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"flags": flags,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to delete a path
//   :values path: A key path
func (config *NsoJsonRpcConfig) Delete(path string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "delete",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"path": path,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get all service points
func (config *NsoJsonRpcConfig) GetServicePoints() (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "get_service_points",
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to get template variables
// This is not xml template variables it is templates in NSO
//   :values name: The name of the template
func (config *NsoJsonRpcConfig) GetTemplateVariables(name string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "get_template_variables",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"name": name,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method for a basic Query in NSO
// This is a convenience method for calling
// start_query, run_query and stop_query This method should not be used for paginated
// results, as it results in performance degradation - use start_query, multiple
// run_query and stop_query instead.
//   :values xpathExpression: A XPATH expression
//   :values resultAs: string, keypath-value, or leaf_value_as_string
func (config *NsoJsonRpcConfig) Query(xpathExpression, resultAs string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "query",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"xpath_expr": xpathExpression,
			"result_as": resultAs,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}


func (config *NsoJsonRpcConfig) StartQuery(xpathExpression, path string, selection []string, chunkSize, initialOffset int, sort []string, sortOrder string, includeTotal bool, contextNode, resultAs string) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "start_query",
		"params": map[string]interface{}{
			"th": config.nsocon.th,
			"xpath_expr": xpathExpression,
			"result_as": resultAs,

		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}
