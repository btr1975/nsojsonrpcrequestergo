package common

import (
	"errors"
	"net"
)

type Common struct {
	protocol, ip, username, password string
	port int
	sslVerify bool
}

// Create a new Common struct
func NewCommon(protocol string, ip string, port int, username string, password string, sslVerify bool) (Common, error)  {

	// Check if protocol is http, or https
	if protocol == "http" || protocol == "https" {
	} else {
		return Common{}, errors.New("Only http, and https is supported!!")
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
			return Common{}, err

		}

	}

	if (port < 1) || (port > 65535) {
		return Common{}, errors.New("valid port range between 1 and 65535")
	}


	return Common{protocol: protocol, ip: ip, port: port, username: username, password: password, sslVerify: sslVerify}, nil

}

