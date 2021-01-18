package nsojsonrpcrequestergo

import (
	"errors"
	"testing"
)

func TestIpV4Address(t *testing.T) {
	scenarios := [] struct{
		input string
		expect string
		rcvError error
	}{
		{input: "192.168.1.1", expect: "192.168.1.1", rcvError: nil},
		{input: "192.168.1.0501", expect: "192.168.1.0501", rcvError: errors.New("not a valid IPv4 address")},
	}

	for _, scenario := range scenarios {
		value, err := IpV4Address(scenario.input)
		if value != scenario.expect {
			t.Errorf("expected %v got %v", scenario.expect, value)
		}

		if err != scenario.rcvError {
			if err.Error() != scenario.rcvError.Error() {
				t.Errorf("expected error %v got %v", scenario.rcvError, err)
			}

		}

	}

}

func TestIpV4MulticastAddress(t *testing.T) {
	scenarios := [] struct{
		input string
		expect string
		rcvError error
	}{
		{input: "192.168.1.1", expect: "192.168.1.1", rcvError: nil},
		{input: "192.168.1.0501", expect: "192.168.1.0501", rcvError: errors.New("not a valid IPv4 address")},
	}

	for _, scenario := range scenarios {
		value, _ := IpV4MulticastAddress(scenario.input)
		if value != scenario.expect {
			t.Errorf("expected %v got %v", scenario.expect, value)
		}

	}
}
