// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"crypto/sha256"
	"fmt"

	"github.com/google/go-configfs-tsm/configfs/configfsi"
	"github.com/google/go-configfs-tsm/configfs/faketsm"
	"github.com/google/go-configfs-tsm/configfs/linuxtsm"
	"github.com/google/go-configfs-tsm/report"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/veraison/cmw"
	"github.com/veraison/da"
)

var (
	createDeviceDirs       []string
	createOutput           string
	createExtendedCCAToken bool
	createDummyTSM         bool
)

var TSMClient configfsi.Client

var createCmd = NewCreateCmd()

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [flags]",
		Short: "TODO",
		Long: `TODO

TODO.

	dat create -d /sys/devices/pci0000:00/0000:00:00.0 -o my-da.cbor
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err error
			)

			if err = checkCreateArgs(args); err != nil {
				return fmt.Errorf("validating arguments: %w", err)
			}

			// initialize TSM client
			if createExtendedCCAToken {
				if !createDummyTSM {
					TSMClient, err = linuxtsm.MakeClient()
					if err != nil {
						return fmt.Errorf("initializing TSM client: %w", err)
					}
				} else {
					// use a dummy TSM client (for testing only)
					TSMClient = &faketsm.Client{
						Subsystems: map[string]configfsi.Client{
							"report": fakeCCATSMClient(),
						},
					}
				}
			}

			// instantiate the DA EAT
			dat := da.NewToken()

			// For each device directory in createDeviceDirs compose create the corresponding device submod
			for _, devDir := range createDeviceDirs {
				if err = addDeviceSubmod(dat, devDir); err != nil {
					return fmt.Errorf("adding device submod for %q: %w", devDir, err)
				}
			}

			// Serialize DA EAT to CBOR
			daBytes, err := dat.ToCBOR()
			if err != nil {
				return fmt.Errorf("serializing DA token to CBOR: %w", err)
			}

			if createExtendedCCAToken {
				// wrap DA EAT in a CCA token
				ccaBytes, err := ratsdCompose(daBytes)
				if err != nil {
					return fmt.Errorf("wrapping DA EAT in the extended CCA token: %w", err)
				}
				daBytes = ccaBytes
			}

			// save to createOutput
			if err = afero.WriteFile(fs, createOutput, daBytes, 0644); err != nil {
				return fmt.Errorf("saving signer DA EAT to file %q: %w", createOutput, err)
			}
			fmt.Printf(">> created %q\n", createOutput)

			return nil
		},
	}

	cmd.Flags().StringVarP(
		&createOutput, "output", "o", "my-dat.cbor", "output file for the DA token (CBOR format)",
	)
	cmd.Flags().StringArrayVarP(
		&createDeviceDirs, "device", "d", []string{}, "device directories to include in the DA token",
	)
	cmd.Flags().BoolVarP(
		&createExtendedCCAToken, "extended-cca-token", "e", false, "wrap the DA EAT in an extended CCA token",
	)
	cmd.Flags().BoolVarP(
		&createDummyTSM, "dummy-tsm", "x", false, "use a dummy TSM client (for testing only)",
	)

	return cmd
}

func ratsdCompose(daBytes []byte) ([]byte, error) {
	nonce := sha256.Sum256(daBytes)
	fmt.Printf(">> generated nonce: %x\n", nonce)

	rsp, err := report.Get(TSMClient, &report.Request{InBlob: nonce[:]})
	if err != nil {
		return nil, fmt.Errorf("getting TSM report: %w", err)
	}

	c, err := cmw.NewCollection("tag:github.com,2025:veraison/ratsd/cmw")
	if err != nil {
		return nil, fmt.Errorf("creating CMW collection: %w", err)
	}

	daNodeMT := `application/eat-ucs+cbor; eat_profile="tag:linaro.org,2025:device#1.0.0"`
	daNode, _ := cmw.NewMonad(daNodeMT, daBytes)
	if err = c.AddCollectionItem("dev", daNode); err != nil {
		return nil, fmt.Errorf("adding trusted device to CMW: %w", err)
	}

	tsmNodeMT := `application/vnd.veraison.tsm-report+cbor`
	tsmNode, _ := cmw.NewMonad(tsmNodeMT, rsp.OutBlob)
	if err = c.AddCollectionItem("tsm", tsmNode); err != nil {
		return nil, fmt.Errorf("adding CCA token to CMW: %w", err)
	}

	return c.MarshalCBOR()
}

func checkCreateArgs(_ []string) error {
	if len(createDeviceDirs) == 0 {
		return fmt.Errorf("no device directories specified")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func addDeviceSubmod(dat *da.Token, devDir string) error {
	fmt.Println(">> adding device submod for", devDir)

	spdmClaims := da.NewSPDMClaims()

	// add certs
	if err := addDeviceCerts(spdmClaims, devDir, fs); err != nil {
		return fmt.Errorf("adding device certs: %w", err)
	}

	// add measurements
	if err := addDeviceMeasurements(spdmClaims, devDir, fs); err != nil {
		return fmt.Errorf("adding device measurements: %w", err)
	}

	if err := dat.AddSubmod(devDir, spdmClaims); err != nil {
		return fmt.Errorf("adding device submod for %q: %w", devDir, err)
	}

	return nil
}

func addDeviceCerts(spdmClaims *da.SPDMClaims, devDir string, fs afero.Fs) error {
	fmt.Println(">> adding device certs for", devDir)

	// read cert from devDir/certificates/slot0
	certPath := fmt.Sprintf("%s/certificates/slot0", devDir)

	certData, err := afero.ReadFile(fs, certPath)
	if err != nil {
		return fmt.Errorf("reading device cert chain from %q: %w", certPath, err)
	}

	return spdmClaims.SetDefaultCert(certData)
}

func addDeviceMeasurements(spdmClaims *da.SPDMClaims, devDir string, fs afero.Fs) error {
	fmt.Println(">> adding device measurements for", devDir)

	for i := 1; i <= 255; i++ {
		typePath := fmt.Sprintf("%s/measurements/type%d", devDir, i)
		digestPath := fmt.Sprintf("%s/measurements/digest%d", devDir, i)
		measPath := fmt.Sprintf("%s/measurements/measurement%d", devDir, i)

		typeData, err := afero.ReadFile(fs, typePath)
		if err != nil {
			// assume no more measurements (should check that err is ENOTFOUND)
			fmt.Println("no more measurements -- stopping at index ", i)
			break
		}

		if len(typeData) != 1 {
			return fmt.Errorf("invalid measurement type length in %q", typePath)
		}

		digestData, err := afero.ReadFile(fs, digestPath)
		if err != nil {
			return fmt.Errorf("reading measurement digest from %q: %w", digestPath, err)
		}

		if len(digestData) != 1 {
			return fmt.Errorf("invalid measurement digest length in %q", digestPath)
		}

		measData, err := afero.ReadFile(fs, measPath)
		if err != nil {
			return fmt.Errorf("reading measurement data from %q: %w", measPath, err)
		}

		meas := da.SPDMMeasurement{
			ComponentType: typeData[0],
			DigestedMeasurement: &da.Digest{
				Algorithm: digestData[0],
				Value:     measData,
			},
		}

		if err := spdmClaims.AddMeasurement(uint8(i), meas); err != nil {
			return fmt.Errorf("adding measurement %d: %w", i, err)
		}

		fmt.Printf(">> added measurement %d (type %d, digest alg %d, len %d)\n", i, typeData[0], digestData[0], len(measData))
	}

	return nil
}
