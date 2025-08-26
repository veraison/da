// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package da

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func TestSPDMClaims_Unmarshal_OK(t *testing.T) {
	tv := []struct {
		file     string
		expected SPDMClaims
	}{
		{"spdm-claims-1.cbor", *exampleSPDMClaims1()},
		{"spdm-claims-2.cbor", *exampleSPDMClaims2()},
	}

	for _, v := range tv {
		b := readTestVectorSlice(t, v.file)
		var actual SPDMClaims
		err := cbor.Unmarshal(b, &actual)
		assert.NoError(t, err)
		assert.Equal(t, v.expected, actual)
	}
}

func TestSPDMClaims_Marshal_OK(t *testing.T) {
	tv := []struct {
		t        SPDMClaims
		expected string
	}{
		{*exampleSPDMClaims1(), "spdm-claims-1.cbor"},
		{*exampleSPDMClaims2(), "spdm-claims-2.cbor"},
	}

	for _, v := range tv {
		actual, err := cbor.Marshal(v.t)
		assert.NoError(t, err)
		assert.Equal(t, readTestVectorSlice(t, v.expected), actual)
	}
}
