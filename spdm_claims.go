// Copyright 2025-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package da

import (
	"crypto/x509"
	"fmt"
)

const SPDMEatProfile = "tag:linaro.org,2025:device-spdm#1.0.0"

type SPDMClaims struct {
	EatProfile            string                          `cbor:"265,keyasint"`
	Measurements          map[uint8]SPDMMeasurement       `cbor:"3802,keyasint,omitempty"`
	Certificates          *SPDMCertificates               `cbor:"3803,keyasint,omitempty"`
	VCA                   *[]byte                         `cbor:"3804,keyasint,omitempty"`
	Challenge             *SPDMMeasurementBlocksSignature `cbor:"3807,keyasint,omitempty"`
	DeviceInterfaceReport *TDISPDeviceInterfaceReport     `cbor:"3808,keyasint,omitempty"`
}

func NewSPDMClaims() *SPDMClaims {
	return &SPDMClaims{
		EatProfile: SPDMEatProfile,
	}
}

func (s *SPDMClaims) AddMeasurement(index uint8, meas SPDMMeasurement) error {
	if s.Measurements == nil {
		s.Measurements = make(map[uint8]SPDMMeasurement)
	}
	s.Measurements[index] = meas
	return nil
}

func (s *SPDMClaims) SetDefaultCert(certChain []byte) error {
	certs, err := x509.ParseCertificates(certChain)
	if err != nil {
		return fmt.Errorf("parsing cert chain from %q: %w", certChain, err)
	}

	if len(certs) == 0 {
		return fmt.Errorf("no valid certificates found in cert chain from %q", certChain)
	}

	if s.Certificates == nil {
		s.Certificates = &SPDMCertificates{}
	}

	s.Certificates.DefaultCertSlot = certChain

	return nil
}

type SPDMMeasurementBlocksSignature struct {
	// TODO: define this
}

type TDISPDeviceInterfaceReport struct {
	// TODO: define this
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
