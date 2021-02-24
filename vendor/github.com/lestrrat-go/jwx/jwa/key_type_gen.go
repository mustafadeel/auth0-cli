// this file was auto-generated by internal/cmd/gentypes/main.go: DO NOT EDIT

package jwa

import (
	"fmt"

	"github.com/pkg/errors"
)

// KeyType represents the key type ("kty") that are supported
type KeyType string

// Supported values for KeyType
const (
	EC             KeyType = "EC"  // Elliptic Curve
	InvalidKeyType KeyType = ""    // Invalid KeyType
	OKP            KeyType = "OKP" // Octet string key pairs
	OctetSeq       KeyType = "oct" // Octet sequence (used to represent symmetric keys)
	RSA            KeyType = "RSA" // RSA
)

var allKeyTypes = []KeyType{
	EC,
	OKP,
	OctetSeq,
	RSA,
}

// KeyTypes returns a list of all available values for KeyType
func KeyTypes() []KeyType {
	return allKeyTypes
}

// Accept is used when conversion from values given by
// outside sources (such as JSON payloads) is required
func (v *KeyType) Accept(value interface{}) error {
	var tmp KeyType
	if x, ok := value.(KeyType); ok {
		tmp = x
	} else {
		var s string
		switch x := value.(type) {
		case fmt.Stringer:
			s = x.String()
		case string:
			s = x
		default:
			return errors.Errorf(`invalid type for jwa.KeyType: %T`, value)
		}
		tmp = KeyType(s)
	}
	switch tmp {
	case EC, OKP, OctetSeq, RSA:
	default:
		return errors.Errorf(`invalid jwa.KeyType value`)
	}

	*v = tmp
	return nil
}

// String returns the string representation of a KeyType
func (v KeyType) String() string {
	return string(v)
}
