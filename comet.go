package nsojsonrpcrequestergo

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
	"math/rand"
)

// NsoJsonRpcComet holds a NSO JSON RPC comet needs
type NsoJsonRpcComet struct {
	nsocon nsoJsonConnection
	cometStarted bool
	cometID string
	handles []string

}


// Constructor for a NsoJsonRpcComet
//   :values nsoJson: A NsoJsonConnection
func NewNsoJsonRpcComet(protocol string, ip string, port int, username string, password string, sslVerify bool) (*NsoJsonRpcComet, error)  {
	cometID := fmt.Sprintf("remote-comet-%d", rand.Intn(65000 - 1 + 1) + 1)

	nsoJson, err := newNsoJsonConnection(protocol, ip, port, username, password, sslVerify)

	if err != nil {
		return &NsoJsonRpcComet{}, err
	}

	return &NsoJsonRpcComet{nsocon: *nsoJson, cometStarted: false, cometID: cometID}, nil

}

func (com *NsoJsonRpcComet) StartComet(username, password string) error {
	err := com.checkCometState(false)

	if err != nil {
		return err
	}

	com.cometStarted = true
	err = com.nsocon.NsoLogin(username, password)

	if err != nil {
		return err
	}

	err = com.nsocon.NewTransaction("read", "private", "", "reuse")

	if err != nil {
		return err
	}

	_, err = com.comet()

	if err != nil {
		return err
	}

	return nil

}


func (com *NsoJsonRpcComet) StopComet() error {
	err := com.checkCometState(true)

	if err != nil {
		return err
	}

	err = com.unsubscribe()

	if err != nil {
		return err
	}

	_, err = com.comet()

	if err != nil {
		return err
	}

	err = com.nsocon.NsoLogout()

	if err != nil {
		return err
	}

	com.cometStarted = false

	return nil

}




func (com *NsoJsonRpcComet) CometPoll() (*req.Resp, error)  {
	response, err := com.comet()

	if err != nil {
		return response, err
	}

	// Need to add something about returning result

	return response, nil


}


func (com *NsoJsonRpcComet) SubscribeChanges(path string) (*req.Resp, error)  {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": com.nsocon.id,
		"method": "subscribe_changes",
		"params": map[string]interface{}{
			"comet_id": com.cometID,
			"path": path,

		},
	}

	response, err := com.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	respData := NewNsoJsonResponse()

	newHandle, err := respData.GetCometHandle(response)

	if err != nil {
		return response, err
	}

	com.handles  = append(com.handles, newHandle)

	response, err = com.startSubscription(newHandle)

	if err != nil {
		return response, err
	}

	return response, nil


}



func (com *NsoJsonRpcComet) SubscribePollLeaf(path string, interval int) (*req.Resp, error)  {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": com.nsocon.id,
		"method": "subscribe_poll_leaf",
		"params": map[string]interface{}{
			"comet_id": com.cometID,
			"path": path,
			"interval": interval,

		},
	}

	response, err := com.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	respData := NewNsoJsonResponse()

	newHandle, err := respData.GetCometHandle(response)

	if err != nil {
		return response, err
	}

	com.handles  = append(com.handles, newHandle)

	response, err = com.startSubscription(newHandle)

	if err != nil {
		return response, err
	}

	return response, nil

}



func (com *NsoJsonRpcComet) SubscribeCDBOper(path string) (*req.Resp, error)  {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": com.nsocon.id,
		"method": "subscribe_cdboper",
		"params": map[string]interface{}{
			"comet_id": com.cometID,
			"path": path,

		},
	}

	response, err := com.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	respData := NewNsoJsonResponse()

	newHandle, err := respData.GetCometHandle(response)

	if err != nil {
		return response, err
	}

	com.handles  = append(com.handles, newHandle)

	response, err = com.startSubscription(newHandle)

	if err != nil {
		return response, err
	}

	return response, nil

}


func (com *NsoJsonRpcComet) SubscribeUpgrade() (*req.Resp, error)  {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": com.nsocon.id,
		"method": "subscribe_upgrade",
		"params": map[string]interface{}{
			"comet_id": com.cometID,

		},
	}

	response, err := com.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	respData := NewNsoJsonResponse()

	newHandle, err := respData.GetCometHandle(response)

	if err != nil {
		return response, err
	}

	com.handles  = append(com.handles, newHandle)

	response, err = com.startSubscription(newHandle)

	if err != nil {
		return response, err
	}

	return response, nil

}


func (com *NsoJsonRpcComet) SubscribeJSONRpcBatch() (*req.Resp, error)  {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": com.nsocon.id,
		"method": "subscribe_jsonrpc_batch",
		"params": map[string]interface{}{
			"comet_id": com.cometID,

		},
	}

	response, err := com.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	respData := NewNsoJsonResponse()

	newHandle, err := respData.GetCometHandle(response)

	if err != nil {
		return response, err
	}

	com.handles  = append(com.handles, newHandle)

	response, err = com.startSubscription(newHandle)

	if err != nil {
		return response, err
	}

	return response, nil

}



func (com *NsoJsonRpcComet) GetSubscriptions() (*req.Resp, error)  {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": com.nsocon.id,
		"method": "get_subscriptions",
	}

	response, err := com.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil

}


func (com *NsoJsonRpcComet) comet() (*req.Resp, error)  {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": com.nsocon.id,
		"method": "comet",
		"params": map[string]interface{}{
			"comet_id": com.cometID,

		},
	}

	response, err := com.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil

}

func (com *NsoJsonRpcComet) startSubscription(handle string) (*req.Resp, error)  {
	param := req.Param{
		"jsonrpc": "2.0",
		"id": com.nsocon.id,
		"method": "start_subscription",
		"params": map[string]interface{}{
			"handle": handle,

		},
	}

	response, err := com.nsocon.sendPost(param)

	if err != nil {
		return response, err
	}

	return response, nil


}

func (com *NsoJsonRpcComet) unsubscribe() error  {
	for _, handle := range com.handles {
		param := req.Param{
			"jsonrpc": "2.0",
			"id": com.nsocon.id,
			"method": "unsubscribe",
			"params": map[string]interface{}{
				"handle": handle,

			},
		}

		_, err := com.nsocon.sendPost(param)

		if err != nil {
			return err
		}

	}

	return nil

}

func (com *NsoJsonRpcComet) checkCometState(wantedState bool) error  {

	if com.cometStarted != wantedState {
		if com.cometStarted == true {
			return errors.New("comet is already running")

		} else {
			return errors.New("comet is not running")

		}
	}

	return nil

}
