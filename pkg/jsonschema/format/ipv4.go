package format

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidIPV4Address = errors.New("invalid IPv4 address")

type IPV4 struct{}

var _ FormatChecker = IPV4{}

func (_ IPV4) CheckFormat(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil
	}

	// net.ParseIP accept 0127.0.0.1 as valid IP address.
	// So here we parse the IP address strictly.
	parts := strings.Split(str, ".")
	if len(parts) != 4 {
		return ErrInvalidIPV4Address
	}

	for _, part := range parts {
		// We first let strconv.ParseInt to accept number in any bases.
		v, err := strconv.ParseInt(part, 0, 64)
		if err != nil {
			return ErrInvalidIPV4Address
		}
		if v < 0 || v > 255 {
			return ErrInvalidIPV4Address
		}
		// And then verify part is in base 10.
		decimal := fmt.Sprintf("%d", v)
		if decimal != part {
			return ErrInvalidIPV4Address
		}
	}

	return nil
}
