package nsojsonrpcrequestergo

import (
	"errors"
	"testing"
)

func TestNewNsoJsonRpcHTTPConnectionGoodParams(t *testing.T) {
	scenarios := [] struct{
		protocol, ip, username, password string
		port int
		sslVerify bool
		headers nsoRequestHeaders
		rcvError error
	}{
		{protocol: "http", ip: "192.168.1.1", port: 8080, username: "admin", password: "admin", sslVerify: false, rcvError: nil},
		{protocol: "https", ip: "192.168.1.1", port: 443, username: "user", password: "pass", sslVerify: true, rcvError: nil},
	}


	for _, scenario := range scenarios {
		tempStruct := &nsoJsonRpcHTTPConnection{
			protocol:  scenario.protocol,
			ip:        scenario.ip,
			username:  scenario.username,
			password:  scenario.password,
			port:      scenario.port,
			sslVerify: scenario.sslVerify,
			headers: nsoRequestHeaders{
				ContentType: "application/json",
				Accept:      "application/json",
			},
		}

		rcvStruct, err := newNsoJsonRpcHTTPConnection(scenario.protocol, scenario.ip, scenario.port, scenario.username, scenario.password, scenario.sslVerify)
		if err != nil {
			t.Fail()
		} else {

			if rcvStruct.protocol != tempStruct.protocol {
				t.Errorf("expected %v got %v", tempStruct.protocol, rcvStruct.protocol)
			}
			if rcvStruct.ip != tempStruct.ip {
				t.Errorf("expected %v got %v", tempStruct.ip, rcvStruct.ip)
			}
			if rcvStruct.port != tempStruct.port {
				t.Errorf("expected %v got %v", tempStruct.port, rcvStruct.port)
			}
			if rcvStruct.username != tempStruct.username {
				t.Errorf("expected %v got %v", tempStruct.username, rcvStruct.username)
			}
			if rcvStruct.password != tempStruct.password {
				t.Errorf("expected %v got %v", tempStruct.username, rcvStruct.username)
			}
			if rcvStruct.sslVerify != tempStruct.sslVerify {
				t.Errorf("expected %v got %v", tempStruct.username, rcvStruct.username)
			}
			if rcvStruct.headers != tempStruct.headers {
				t.Errorf("expected %v got %v", tempStruct.headers, rcvStruct.headers)
			}
		}

	}

}

func TestNewNsoJsonRpcHTTPConnectionBadParams(t *testing.T) {
	scenarios := [] struct{
		protocol, ip, username, password string
		port int
		sslVerify bool
		headers nsoRequestHeaders
		rcvError error
	}{
		{protocol: "ssh", ip: "192.168.1.1", port: 8080, username: "admin", password: "admin", sslVerify: false, rcvError: errors.New("only http, and https is supported")},
		{protocol: "https", ip: "192.168.1.1000", port: 443, username: "user", password: "pass", sslVerify: true, rcvError: errors.New("not a valid IPv4 address")},
		{protocol: "https", ip: "192.168.1.1", port: 0, username: "user", password: "pass", sslVerify: true, rcvError: errors.New("valid port range between 1 and 65535")},
		{protocol: "https", ip: "192.168.1.1", port: 65536, username: "user", password: "pass", sslVerify: true, rcvError: errors.New("valid port range between 1 and 65535")},
	}

	for _, scenario := range scenarios {
		_, err := newNsoJsonRpcHTTPConnection(scenario.protocol, scenario.ip, scenario.port, scenario.username, scenario.password, scenario.sslVerify)
		if err == nil {
			t.Fail()

		} else {
			if err.Error() != scenario.rcvError.Error() {
				t.Errorf("expected error %v got %v", scenario.rcvError, err)
			}
		}


	}

}

func TestNsoJsonRpcHTTPConnection_NsoUrl(t *testing.T) {
	scenarios := [] struct{
		protocol, ip, username, password string
		port int
		sslVerify bool
		headers nsoRequestHeaders
		expect string
	}{
		{protocol: "http", ip: "192.168.1.1", port: 8080, username: "admin", password: "admin", sslVerify: false, expect: "http://192.168.1.1:8080/jsonrpc"},
		{protocol: "https", ip: "192.168.1.1", port: 443, username: "user", password: "pass", sslVerify: true, expect: "https://192.168.1.1:443/jsonrpc"},
	}

	for _, scenario := range scenarios {
		rcvStruct, err := newNsoJsonRpcHTTPConnection(scenario.protocol, scenario.ip, scenario.port, scenario.username, scenario.password, scenario.sslVerify)
		if err != nil {
			t.Fail()
		} else {
			if rcvStruct.NsoUrl() != scenario.expect {
				t.Errorf("expected %v got %v", scenario.expect, rcvStruct.NsoUrl())
			}

		}

	}

}
