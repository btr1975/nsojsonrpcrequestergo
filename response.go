package nsojsonrpcrequestergo

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
	"strings"
)

// NsoJsonResponse holds a NSO JSON RPC Response
// The tags help to convert fields to lowercase
type NsoJsonResponse struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Result  map[string]interface{} `json:"result"`
	ID      int                    `json:"id"`
	Error   map[string]interface{} `json:"error"`
}

// Constructor to create a new NsoJsonResponse struct
func NewNsoJsonResponse() *NsoJsonResponse {

	return &NsoJsonResponse{}

}

// Method to convert JSON Response Body to a map
//   :values response: *req.Resp
func (r *NsoJsonResponse) ResponseToStruct(response *req.Resp) (*NsoJsonResponse, error) {
	err := response.ToJSON(&r)

	if err != nil {
		return r, err
	}

	return r, nil

}

// Method to get the transaction handle
//   :values response: *req.Resp
func (r *NsoJsonResponse) GetTransactionHandle(response *req.Resp) float64 {
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

// Method to get the query handle
//   :values response: *req.Resp
func (r *NsoJsonResponse) GetQueryHandle(response *req.Resp) float64 {
	var qh float64

	_ = response.ToJSON(&r)

	for key, value := range r.Result {
		if key == "qh" {
			qh = value.(float64)
			break
		}
	}

	return qh

}

// Method to convert the query results string to a array
//   :values result: result string from query
func (r *NsoJsonResponse) fixQueryResults(result string) []string {
	remLeftDblBracket := strings.Replace(result, "[[", "", -1)
	remRightDblBracket := strings.Replace(remLeftDblBracket, "]]", "", -1)
	csvFmt := strings.Replace(remRightDblBracket, "] [", ",", -1)
	return strings.Split(csvFmt, ",")

}

// Method to get the query results
//   :values response: *req.Resp
func (r *NsoJsonResponse) GetQueryResults(response *req.Resp) ([]string, error) {
	_, err := r.ResponseToStruct(response)
	if err != nil {
		return []string{}, err
	}

	resutlData := r.Result
	for k, v := range resutlData {
		if k == "results" {
			return r.fixQueryResults(fmt.Sprintf("%s", v)), nil
		}
	}

	return []string{}, errors.New("could not find results")

}

// Method to get the comet handle
//   :values response: *req.Resp
func (r *NsoJsonResponse) GetCometHandle(response *req.Resp) (string, error) {
	_, err := r.ResponseToStruct(response)

	if err != nil {
		return "", err
	}

	resutlData := r.Result
	for k, v := range resutlData {
		if k == "handle" {
			return fmt.Sprintf("%s", v), nil
		}
	}

	return "", errors.New("could not find handle")

}
