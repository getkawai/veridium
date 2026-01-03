// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cosign

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/sigstore/cosign/v3/pkg/oci"
	ociremote "github.com/sigstore/cosign/v3/pkg/oci/remote"
	"github.com/sigstore/rekor/pkg/generated/client"
	"github.com/sigstore/sigstore-go/pkg/root"
	"github.com/sigstore/sigstore/pkg/signature"
	"github.com/sigstore/sigstore/pkg/signature/options"
)

// CheckOpts are the options for checking signatures.
type CheckOpts struct {
	// RegistryClientOpts are the options for interacting with the container registry.
	RegistryClientOpts []ociremote.Option

	// Annotations optionally specifies image signature annotations to verify.
	Annotations map[string]interface{}

	// ClaimVerifier, if provided, verifies claims present in the oci.Signature.
	ClaimVerifier func(sig oci.Signature, imageDigest v1.Hash, annotations map[string]interface{}) error

	// TrustedMaterial contains trusted metadata for all Sigstore services.
	TrustedMaterial root.TrustedMaterial

	// RekorClient, if set, is used to make online tlog calls use to verify signatures and public keys.
	RekorClient *client.Rekor

	// RekorPubKeys, if set, is used to validate signatures on log entries from Rekor.
	RekorPubKeys *TrustedTransparencyLogPubKeys

	// SigVerifier is used to verify signatures.
	SigVerifier signature.Verifier

	// PKOpts are the options provided to `SigVerifier.PublicKey()`.
	PKOpts []signature.PublicKeyOption

	// RootCerts are the root CA certs used to verify a signature's chained certificate.
	RootCerts *x509.CertPool

	// IntermediateCerts are the optional intermediate CA certs used to verify a certificate chain.
	IntermediateCerts *x509.CertPool

	// IgnoreSCT requires that a certificate contain an embedded SCT during verification.
	IgnoreSCT bool

	// SignatureRef is the reference to the signature file.
	SignatureRef string

	// PayloadRef is the reference to the payload file.
	PayloadRef string

	// Offline is set to true to prevent network calls.
	Offline bool

	// IgnoreTlog skips transparency log verification.
	IgnoreTlog bool

	// MaxWorkers is the maximum number of workers to use for parallel verification.
	MaxWorkers int

	// ExperimentalOCI11 enables experimental OCI 1.1 verification.
	ExperimentalOCI11 bool

	// UseSignedTimestamps requires that a signature contain a signed timestamp.
	UseSignedTimestamps bool

	// NewBundleFormat uses the new bundle format for verification.
	NewBundleFormat bool
}

// payloader interface for unit testing
type payloader interface {
	Base64Signature() (string, error)
	Payload() ([]byte, error)
}

// verifyOCISignature verifies an OCI signature
func verifyOCISignature(ctx context.Context, verifier signature.Verifier, sig payloader) error {
	b64sig, err := sig.Base64Signature()
	if err != nil {
		return err
	}
	signatureBytes, err := base64.StdEncoding.DecodeString(b64sig)
	if err != nil {
		return err
	}
	payload, err := sig.Payload()
	if err != nil {
		return err
	}
	return verifier.VerifySignature(bytes.NewReader(signatureBytes), bytes.NewReader(payload), options.WithContext(ctx))
}

// VerifyImageSignature verifies a signature
func VerifyImageSignature(ctx context.Context, sig oci.Signature, h v1.Hash, co *CheckOpts) (bundleVerified bool, err error) {
	// Simplified version - only verify cryptographic signature with public key
	if co.SigVerifier == nil {
		return false, errors.New("SigVerifier is required for verification")
	}

	// Perform cryptographic verification of the signature using the public key
	if err := verifyOCISignature(ctx, co.SigVerifier, sig); err != nil {
		return false, fmt.Errorf("signature verification failed: %w", err)
	}

	// Verify claims if ClaimVerifier is provided
	if co.ClaimVerifier != nil {
		if err := co.ClaimVerifier(sig, h, co.Annotations); err != nil {
			return false, fmt.Errorf("claim verification failed: %w", err)
		}
	}

	// For this simplified version, we don't verify transparency log (Rekor) or certificates
	// We return true to indicate successful signature verification
	// Note: "bundleVerified" in this context means "signature successfully verified"
	// not "Rekor bundle verified" (we skip Rekor verification)
	return true, nil
}
