// Package dmapi provides the DistroMate DM API SDK for Go.
package dmapi

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// DmApi is the high-level SDK client.
type DmApi struct {
	pipeTimeout uint32
}

// ShouldSkipCheck validates local dev license and reports whether runtime checks can be skipped.
func ShouldSkipCheck(appID string, publicKey string) (bool, error) {
	if os.Getenv("DM_PIPE") != "" && os.Getenv("DM_API_PATH") != "" {
		return false, nil
	}

	resolvedAppID := strings.TrimSpace(appID)
	resolvedPublicKey := strings.TrimSpace(publicKey)
	if resolvedAppID == "" {
		resolvedAppID = strings.TrimSpace(os.Getenv("DM_APP_ID"))
	}
	if resolvedPublicKey == "" {
		resolvedPublicKey = strings.TrimSpace(os.Getenv("DM_PUBLIC_KEY"))
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
// pipeTimeout is optional and defaults to 5000ms.
func New(dllPath string, pipeTimeout ...uint32) (*DmApi, error) {
	if err := ensureDLL(dllPath); err != nil {
		return nil, err
	}

	timeout := uint32(defaultPipeTimeout)
	if len(pipeTimeout) > 0 {
		timeout = pipeTimeout[0]
	}

	return &DmApi{pipeTimeout: timeout}, nil
}

// SetPipeTimeout updates the timeout used by update IPC methods.
func (api *DmApi) SetPipeTimeout(timeout uint32) {
	api.pipeTimeout = timeout
}

func (api *DmApi) resolvePipeTimeout() uint32 {
	if api.pipeTimeout == 0 {
		return defaultPipeTimeout
	}
	return api.pipeTimeout
}

func (api *DmApi) connectPipe() bool {
	pipe := strings.TrimSpace(os.Getenv("DM_PIPE"))
	if pipe == "" {
		return false
	}

	pipeBytes := cStringBytes(pipe)
	ret, _, _ := procConnect.Call(
		uintptr(unsafe.Pointer(&pipeBytes[0])),
		uintptr(api.resolvePipeTimeout()),
	)
	return int32(ret) == 0
}

func (api *DmApi) callPipeJSON(proc *syscall.Proc, args ...uintptr) map[string]interface{} {
	if !api.connectPipe() {
		return nil
	}
	defer procClose.Call()

	ptr, _, _ := proc.Call(args...)
	if ptr == 0 {
		return nil
	}

	return parseResponseData(ptrToOwnedString(ptr))
}

func (api *DmApi) callPipeAccepted(proc *syscall.Proc, args ...uintptr) bool {
	if !api.connectPipe() {
		return false
	}
	defer procClose.Call()

	ret, _, _ := proc.Call(args...)
	return int32(ret) == 1
}
