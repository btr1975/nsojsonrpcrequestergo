package nsojsonrpcrequestergo

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
)

// NsoJsonRpcConfig holds a NSO JSON RPC config needs
type NsoJsonRpcConfig struct {
	nsocon nsoJsonConnection

}

// Constructor for a NsoJsonRpcConfig
//   :values protocol: http, https
//   :values ip: a IPv4 address, or a CNAME
//   :values port: 1 to 65535
//   :values username: A username
//   :values password: A password
//   :values sslVerify: true to verify SSL, false not to
func NewNsoJsonRpcConfig(protocol string, ip string, port int, username string, password string, sslVerify bool) (*NsoJsonRpcConfig, error)  {

	nsoJson, err := newNsoJsonConnection(protocol, ip, port, username, password, sslVerify)

	if err != nil {
		return &NsoJsonRpcConfig{}, err
	}

	return &NsoJsonRpcConfig{nsocon: *nsoJson}, nil

}

// Method to login to the NSO Server
func (config *NsoJsonRpcConfig) NsoLogin() error {
	err := config.nsocon.NsoLogin()

	if err != nil {
		return err
	}

	return nil

}

// Method to logout to the NSO Server
func (config *NsoJsonRpcConfig) NsoLogout() error {
	err := config.nsocon.NsoLogout()

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
func (config *NsoJsonRpcConfig) NewTransaction(mode, confMode, tag, onPendingChanges string) error {
	err := config.nsocon.NewTransaction(mode, confMode, tag, onPendingChanges)

	if err != nil {
		return err
	}

	return nil
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

// QueryObject contains a compext query structure
type QueryObject struct {
	xpathExpression, path, sortOrder, contextNode, resultAs  string
	selection, sort []string
	chunkSize, initialOffset int
	includeTotal bool
	qh float64
}

// Constructor for a QueryObject
// If xpathExpression is defined path will be ignored
//   :values xpathExpression: A XPATH expression or leave blank to use a keypath instead
//   :values path: A keypath epression
//   :values selection: An array of leaf selections use empty array at your own risk
//   :values chunkSize: If set to 0 all data is returned, any other value will break data into chunks
//   :values initialOffset: If set to 0 thing is done any other value sets the offset
//   :values sort: Array of XPATH expressions use a blank array to not use
//   :values sortOrder: ascending, or descending "" to not use
//   :values includeTotal: true to include total records, false to not
//   :values contextNode: A keypath optional use "" to not use
//   :values resultAs: string, keypath-value, or leaf_value_as_string
func NewQueryObject (xpathExpression, path string, selection []string, chunkSize, initialOffset int, sort []string, sortOrder string, includeTotal bool, contextNode, resultAs string) (*QueryObject, error) {
	var expression, usepath string

	if xpathExpression != "" {
		expression = xpathExpression

	} else if path != "" {
		usepath = path

	} else {
		return &QueryObject{}, errors.New("either xpathExpression needs to be given or path")

	}

	return &QueryObject{xpathExpression: expression, path: usepath, selection: selection, chunkSize: chunkSize, initialOffset: initialOffset, sort: sort, sortOrder: sortOrder, includeTotal: includeTotal, contextNode: contextNode, resultAs: resultAs}, nil
}

// Method to start a complex query
//   :vaules QueryObject: A QueryObject
func (config *NsoJsonRpcConfig) StartQuery(queryObject *QueryObject) error {
	params := map[string]interface{}{
		"th": config.nsocon.th,
	}
	if queryObject.xpathExpression != "" {
		params["xpath_expr"] = queryObject.xpathExpression
		if len(queryObject.selection) > 0 {
			params["selection"] = queryObject.selection
		}

		if len(queryObject.sort) > 0 {
			params["sort"] = queryObject.sort
		}

	} else {
		params["path"] = queryObject.path
		if queryObject.contextNode != "" {
			params["context_node"] = queryObject.contextNode
		}

	}
	params["chunk_size"] = queryObject.chunkSize
	params["initial_offset"] = queryObject.initialOffset
	if queryObject.sortOrder != "" {
		params["sort_order"] = queryObject.sortOrder
	}

	params["include_total"] = queryObject.includeTotal
	params["result_as"] = queryObject.resultAs

	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "start_query",
		"params": params,
	}

	response, err := config.nsocon.sendPost(param)

	nsoResponse := NewNsoJsonResponse()
	queryObject.qh = nsoResponse.GetQueryHandle(response)

	if err != nil {
		return err
	}

	return nil
}

// Method to run a complex query
//   :values queryHandle: A Query Handle this comes from using the StartQuery method
func (config *NsoJsonRpcConfig) RunQuery(queryObject *QueryObject) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "run_query",
		"params": map[string]interface{}{
			"qh": queryObject.qh,
		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}

// Method to reset a complex query
//   :values queryHandle: A Query Handle this comes from using the StartQuery method
func (config *NsoJsonRpcConfig) ResetQuery(queryObject *QueryObject) (*req.Resp, error) {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "reset_query",
		"params": map[string]interface{}{
			"qh": queryObject.qh,
		},
	}

	response, err := config.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil
}


// Method to stop a complex query
//   :values queryHandle: A Query Handle this comes from using the StartQuery method
func (config *NsoJsonRpcConfig) StopQuery(queryObject *QueryObject) error {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": config.nsocon.id,
		"method": "stop_query",
		"params": map[string]interface{}{
			"qh": queryObject.qh,
		},
	}

	_, err := config.nsocon.sendPost(param)

	if err != nil {
		return err
	}

	return nil
}
