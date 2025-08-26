// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package da

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func TestDAToken_Unmarshal_OK(t *testing.T) {
	tv := []struct {
		file     string
		expected DAToken
	}{
		{"eat-da-1.cbor", *exampleDAToken1()},
	}

	for _, v := range tv {
		b := readTestVectorSlice(t, v.file)
		var actual DAToken
		err := cbor.Unmarshal(b, &actual)
		assert.NoError(t, err)
		assert.Equal(t, v.expected, actual)
	}
}
