// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package da

type SPDMClaims struct {
	EatProfile   string                          `cbor:"265,keyasint"`
	Measurements map[uint8]SPDMMeasurement       `cbor:"1,keyasint"`
	Signature    *SPDMMeasurementBlocksSignature `cbor:"signature,omitempty"`
	Certificates SPDMCertificates                `cbor:"2,keyasint"`
}

func NewSPDMClaims() *SPDMClaims {
	return &SPDMClaims{
		EatProfile: "tag:linaro.org,2025:device-spdm#1.0.0",
	}
}

type SPDMMeasurementBlocksSignature struct {
	// TODO
}

type ComponentType int

const (
	ComponentTypeUnknown ComponentType = iota
	ComponentTypeImmutableROM
	ComponentTypeMutableFirmware
	ComponentTypeHardwareConfig
	ComponentTypeFirmwareConfig
	ComponentTypeFreeformMeasurementManifest
	ComponentTypeDeviceMode
	ComponentTypeMutableFirmwareVersion
	ComponentTypeMutableFirmwareSVN
	ComponentTypeHashExtendMeasurement
	ComponentTypeInformational
	ComponentTypeStructuredMeasurementManifest
)

type SPDMMeasurement struct {
	ComponentType       uint8   `cbor:"1,keyasint"`
	DigestedMeasurement *Digest `cbor:"2,keyasint,omitempty"`
	RawMeasurement      *[]byte `cbor:"3,keyasint,omitempty"`
	Signature           *[]byte `cbor:"4,keyasint,omitempty"`
}

type Digest struct {
	_         struct{} `cbor:",toarray"`
	Algorithm uint8    // TODO(tho) add text as allowed type (or reuse corim's)
	Value     []byte
}

type SPDMCertificates struct {
	DefaultCertSlot []byte  `cbor:"0,keyasint"`
	AuxCertSlot1    *[]byte `cbor:"1,keyasint,omitempty"`
	AuxCertSlot2    *[]byte `cbor:"2,keyasint,omitempty"`
	AuxCertSlot3    *[]byte `cbor:"3,keyasint,omitempty"`
	AuxCertSlot4    *[]byte `cbor:"4,keyasint,omitempty"`
	AuxCertSlot5    *[]byte `cbor:"5,keyasint,omitempty"`
	AuxCertSlot6    *[]byte `cbor:"6,keyasint,omitempty"`
	AuxCertSlot7    *[]byte `cbor:"7,keyasint,omitempty"`
}
