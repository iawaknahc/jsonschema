package format

import (
	"errors"
	"net"
)

var ErrInvalidIPV4Address = errors.New("invalid IPv4 address")

type IPV4 struct{}

var _ FormatChecker = IPV4{}

func (_ IPV4) CheckFormat(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil
	}
	ip := net.ParseIP(str)
	if ip == nil {
		return ErrInvalidIPV4Address
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return ErrInvalidIPV4Address
	}
	return nil
}
