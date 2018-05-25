package dependency

import (
	"crypto/x509"
	"encoding/pem"
	"time"
)

var (
	// VaultDefaultLeaseDuration is the default lease duration in seconds.
	VaultDefaultLeaseDuration = 5 * time.Minute
)

// Secret is a vault secret.
type Secret struct {
	RequestID     string
	LeaseID       string
	LeaseDuration int
	Renewable     bool

	// Data is the actual contents of the secret. The format of the data
	// is arbitrary and up to the secret backend.
	Data map[string]interface{}
}

// leaseDurationOrDefault returns a value or the default lease duration.
func leaseDurationOrDefault(d int) int {
	if d == 0 {
		return int(VaultDefaultLeaseDuration.Nanoseconds() / 1000000000)
	}
	return d
}

// vaultRenewDuration accepts a given renew duration (lease duration) and
// returns the cooresponding time.Duration. If the duration is 0 (not provided),
// this falls back to the VaultDefaultLeaseDuration.
func vaultRenewDuration(d int) time.Duration {
	dur := time.Duration(d/2.0) * time.Second
	if dur == 0 {
		dur = VaultDefaultLeaseDuration
	}
	return dur
}

// durationFromCert gets a "lease" duration in seconds from cert data
func durationFromCert(certData string) int {
	block, _ := pem.Decode([]byte(certData))
	if block == nil {
		return -1
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return -1
	}

	return int(cert.NotAfter.Sub(cert.NotBefore).Seconds())
}
