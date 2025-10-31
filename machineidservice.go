package main

import (
	"github.com/kawai-network/veridium/machineid"
)

// MachineIDService provides access to machine identification functionality
type MachineIDService struct{}

// GetID returns the platform specific machine ID of the current host OS
func (m *MachineIDService) GetID() (string, error) {
	return machineid.ID()
}

// GetProtectedID returns a hashed version of the machine ID in a cryptographically secure way,
// using a fixed, application-specific key.
// Internally, this function calculates HMAC-SHA256 of the application ID, keyed by the machine ID.
func (m *MachineIDService) GetProtectedID(appID string) (string, error) {
	return machineid.ProtectedID(appID)
}
