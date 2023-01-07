package nsojsonrpcrequestergo

import (
	"errors"
	"net"
)

// IpV4Address verifies if a given string a IPv4 Address
//
//	:values address: The address to verify
func IpV4Address(address string) (string, error) {
	// Convert ip to a IPv4 Address if possible
	ipv4address := net.ParseIP(address)

	if ipv4address == nil {
		return address, errors.New("not a valid IPv4 address")

	} else {

		return ipv4address.String(), nil
	}

}

// IpV4UnicastAddress verifies if a net.IP is a IPv4 Unicast address
//
//	:values address: The address to verify
func IpV4UnicastAddress(address string) (string, error) {
	_, err := IpV4Address(address)

	if err != nil {
		return address, err

	}

	isUcast := net.IP.IsGlobalUnicast(net.ParseIP(address))

	if isUcast != true {
		return address, errors.New("not a valid IPv4 unicast address")

	}

	return address, nil

}

// IpV4MulticastAddress verifies if a net.IP is a IPv4 Multicast address
//
//	:values address: The address to verify
func IpV4MulticastAddress(address string) (string, error) {
	_, err := IpV4Address(address)

	if err != nil {
		return address, err

	}

	isMcast := net.IP.IsMulticast(net.ParseIP(address))

	if isMcast != true {
		return address, errors.New("not a valid IPv4 multicast address")

	}

	return address, nil

}
