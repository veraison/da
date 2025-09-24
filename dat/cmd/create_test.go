// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"testing"

	"github.com/google/go-configfs-tsm/configfs/configfsi"
	"github.com/google/go-configfs-tsm/configfs/faketsm"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_CreateCmd_ok(t *testing.T) {
	cmd := NewCreateCmd()

	testFS1(t)

	args := []string{
		"--device=/sys/devices/pci0000:00/0000:00:00.0",
		"--output=./test-dat.cbor",
	}
	cmd.SetArgs(args)

	err := cmd.Execute()
	assert.NoError(t, err)

	actual, err := afero.ReadFile(fs, "test-dat.cbor")
	assert.NoError(t, err)

	expected := testEATDA

	assert.Equal(t, expected, actual)
}

func Test_CreateCmd_extended(t *testing.T) {
	cmd := NewCreateCmd()

	testFS1(t)

	args := []string{
		"--device=/sys/devices/pci0000:00/0000:00:00.0",
		"--output=./test-dat.cbor",
		"--extended-cca-token",
		"--dummy-tsm",
	}
	cmd.SetArgs(args)

	// provide a custom TSMClient that returns a fixed CCA token
	TSMClient = &faketsm.Client{
		Subsystems: map[string]configfsi.Client{
			"report": fakeCCATSMClient(),
		},
	}

	err := cmd.Execute()
	assert.NoError(t, err)

	actual, err := afero.ReadFile(fs, "test-dat.cbor")
	assert.NoError(t, err)

	expected := testTSMReport

	assert.Equal(t, expected, actual)
}
