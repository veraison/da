// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package da

import (
	"github.com/fxamacker/cbor/v2"
)

type DAToken struct {
	EatProfile string                `cbor:"265,keyasint"`
	EatNonce   [64]byte              `cbor:"10,keyasint"`
	EatSubmods map[string]SPDMClaims `cbor:"266,keyasint"`
}

// NewDAToken creates a new DAToken instance.
func NewDAToken() *DAToken {
	return &DAToken{
		EatProfile: "tag:linaro.org,2025:device#1.0.0",
	}
}

func (d *DAToken) AddSubmod(name string, claims *SPDMClaims) error {
	if d.EatSubmods == nil {
		d.EatSubmods = make(map[string]SPDMClaims)
	}
	d.EatSubmods[name] = *claims
	return nil
}

func (d *DAToken) ToCBOR() ([]byte, error) {
	em, err := cbor.CoreDetEncOptions().EncMode()
	if err != nil {
		return nil, err
	}
	return em.Marshal(d)
}

func (d *DAToken) FromCBOR(data []byte) error {
	return cbor.Unmarshal(data, d)
}
