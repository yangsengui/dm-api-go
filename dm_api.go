// Package dmapi provides the DistroMate DM API SDK for Go.
package dmapi

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// DmApi is the high-level SDK client.
type DmApi struct {
	licenseCallbackFn   func()
	licenseCallbackAddr uintptr
}

// ShouldSkipCheck validates local dev license and reports whether runtime checks can be skipped.
func ShouldSkipCheck(appID string, publicKey string) (bool, error) {
	if os.Getenv(envDmLauncherEndpoint) != "" && os.Getenv(envDmLauncherToken) != "" {
		return false, nil
	}

	resolvedAppID := strings.TrimSpace(appID)
	resolvedPublicKey := strings.TrimSpace(publicKey)
	if resolvedAppID == "" {
		resolvedAppID = strings.TrimSpace(os.Getenv(envDmAppID))
	}
	if resolvedPublicKey == "" {
		resolvedPublicKey = strings.TrimSpace(os.Getenv(envDmPublicKey))
	}

	if resolvedAppID == "" || resolvedPublicKey == "" {
		return false, errors.New("app identity is required for dev-license checks. Provide appID/publicKey or set DM_APP_ID and DM_PUBLIC_KEY")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return false, errors.New(devLicenseErrorText)
	}

	pubkeyPath := filepath.Join(home, ".distromate-cli", "dev_licenses", resolvedAppID, "pubkey")
	raw, err := os.ReadFile(pubkeyPath)
	if err != nil {
		return false, errors.New(devLicenseErrorText)
	}

	devPubKey := strings.TrimSpace(string(raw))
	if devPubKey == "" || devPubKey != resolvedPublicKey {
		return false, errors.New(devLicenseErrorText)
	}

	return true, nil
}

// New initializes SDK and loads dm_api.dll.
func New(dllPath string) (*DmApi, error) {
	if err := ensureDLL(dllPath); err != nil {
		return nil, err
	}

	return &DmApi{}, nil
}

// GetActivationErrorName resolves activation error code to symbolic name.
func (api *DmApi) GetActivationErrorName(code uint32) string {
	if name, ok := activationErrorNames[code]; ok {
		return name
	}
	return "UNKNOWN(" + strconvU32(code) + ")"
}

func strconvU32(value uint32) string {
	const digits = "0123456789"
	if value == 0 {
		return "0"
	}

	buf := [10]byte{}
	i := len(buf)
	for value > 0 {
		i--
		buf[i] = digits[value%10]
		value /= 10
	}
	return string(buf[i:])
}
