package format

import (
	"errors"
	"fmt"
	"net/mail"
)

var ErrEmailAddressWithName = errors.New("invalid email address")

type Email struct{}

var _ FormatChecker = Email{}

func (Email) CheckFormat(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return nil
	}

	addr, err := mail.ParseAddress(str)
	if err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	if addr.Name != "" {
		return ErrEmailAddressWithName
	}

	return nil
}
