// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package da

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
