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
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/sigstore/cosign/v3/pkg/cosign/env"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
	"github.com/sigstore/sigstore/pkg/tuf"
)

// TransparencyLogPubKey contains the ECDSA verification key and the current status
// of the key according to TUF metadata, whether it's active or expired.
type TransparencyLogPubKey struct {
	PubKey crypto.PublicKey
	Status tuf.StatusKind
}

// TrustedTransparencyLogPubKeys is a map of TransparencyLog public keys indexed by log ID
// that's used in verification.
type TrustedTransparencyLogPubKeys struct {
	// A map of keys indexed by log ID
	Keys map[string]TransparencyLogPubKey
}

// NewTrustedTransparencyLogPubKeys creates a new TrustedTransparencyLogPubKeys
func NewTrustedTransparencyLogPubKeys() TrustedTransparencyLogPubKeys {
	return TrustedTransparencyLogPubKeys{Keys: make(map[string]TransparencyLogPubKey, 0)}
}

// GetTransparencyLogID generates a SHA256 hash of a DER-encoded public key.
// (see RFC 6962 S3.2)
// In CT V1 the log id is a hash of the public key.
func GetTransparencyLogID(pub crypto.PublicKey) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(pubBytes)
	return hex.EncodeToString(digest[:]), nil
}

// AddTransparencyLogPubKey returns a log ID and RekorPubKey from a given
// byte-array representing the PEM-encoded Rekor key and a status.
func (t *TrustedTransparencyLogPubKeys) AddTransparencyLogPubKey(pemBytes []byte, status tuf.StatusKind) error {
	pubKey, err := cryptoutils.UnmarshalPEMToPublicKey(pemBytes)
	if err != nil {
		return err
	}
	keyID, err := GetTransparencyLogID(pubKey)
	if err != nil {
		return err
	}
	t.Keys[keyID] = TransparencyLogPubKey{PubKey: pubKey, Status: status}
	return nil
}

// This is the rekor transparency log public key target name
var rekorTargetStr = `rekor.pub`

// GetRekorPubs retrieves trusted Rekor public keys from the embedded or cached
// TUF root. If expired, makes a network call to retrieve the updated targets.
// There are two Env variable that can be used to override this behaviour:
// SIGSTORE_REKOR_PUBLIC_KEY - If specified, location of the file that contains
// the Rekor Public Key on local filesystem
func GetRekorPubs(ctx context.Context) (*TrustedTransparencyLogPubKeys, error) {
	publicKeys := NewTrustedTransparencyLogPubKeys()
	altRekorPub := env.Getenv(env.VariableSigstoreRekorPublicKey)

	if altRekorPub != "" {
		raw, err := os.ReadFile(altRekorPub)
		if err != nil {
			return nil, fmt.Errorf("error reading alternate Rekor public key file: %w", err)
		}
		if err := publicKeys.AddTransparencyLogPubKey(raw, tuf.Active); err != nil {
			return nil, fmt.Errorf("AddRekorPubKey: %w", err)
		}
	} else {
		tufClient, err := tuf.NewFromEnv(ctx)
		if err != nil {
			return nil, err
		}
		targets, err := tufClient.GetTargetsByMeta(tuf.Rekor, []string{rekorTargetStr})
		if err != nil {
			return nil, err
		}
		for _, t := range targets {
			if err := publicKeys.AddTransparencyLogPubKey(t.Target, t.Status); err != nil {
				return nil, fmt.Errorf("AddRekorPubKey: %w", err)
			}
		}
	}

	if len(publicKeys.Keys) == 0 {
		return nil, errors.New("none of the Rekor public keys have been found")
	}

	return &publicKeys, nil
}
