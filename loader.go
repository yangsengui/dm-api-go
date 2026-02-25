package dmapi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

const (
	defaultDLLName      = "dm_api.dll"
	defaultPipeTimeout  = 5000
	defaultBufferSize   = 256
	defaultVersionSize  = 32
	defaultModeBufSize  = 64
	devLicenseErrorText = "Development license is missing or corrupted. Run `distromate sdk renew` to regenerate the dev certificate."
)

var (
	dll      *syscall.DLL
	loadErr  error
	loadOnce sync.Once

	procConnect                        *syscall.Proc
	procClose                          *syscall.Proc
	procGetVersion                     *syscall.Proc
	procRestart                        *syscall.Proc
	procLastError                      *syscall.Proc
	procCheckForUpdates                *syscall.Proc
	procDownloadUpdate                 *syscall.Proc
	procGetUpdateState                 *syscall.Proc
	procWaitUpdateState                *syscall.Proc
	procQuitAndInstall                 *syscall.Proc
	procJsonToCanonical                *syscall.Proc
	procFreeString                     *syscall.Proc
	procSetProductData                 *syscall.Proc
	procSetProductId                   *syscall.Proc
	procSetDataDirectory               *syscall.Proc
	procSetDebugMode                   *syscall.Proc
	procSetCustomDeviceFingerprint     *syscall.Proc
	procSetLicenseKey                  *syscall.Proc
	procSetActivationMetadata          *syscall.Proc
	procActivateLicense                *syscall.Proc
	procActivateLicenseOffline         *syscall.Proc
	procGenerateOfflineDeactivationReq *syscall.Proc
	procGetLastActivationError         *syscall.Proc
	procIsLicenseGenuine               *syscall.Proc
	procIsLicenseValid                 *syscall.Proc
	procGetServerSyncGracePeriodExpiry *syscall.Proc
	procGetActivationMode              *syscall.Proc
	procGetLicenseKey                  *syscall.Proc
	procGetLicenseExpiryDate           *syscall.Proc
	procGetLicenseCreationDate         *syscall.Proc
	procGetLicenseActivationDate       *syscall.Proc
	procGetActivationCreationDate      *syscall.Proc
	procGetActivationLastSyncedDate    *syscall.Proc
	procGetActivationId                *syscall.Proc
	procGetLibraryVersion              *syscall.Proc
	procReset                          *syscall.Proc
)

func resolveDLLPath(dllPath string) string {
	resolved := strings.TrimSpace(dllPath)
	if resolved == "" {
		resolved = strings.TrimSpace(os.Getenv("DM_API_PATH"))
	}
	if resolved == "" {
		resolved = defaultDLLName
	}

	if filepath.IsAbs(resolved) {
		return resolved
	}

	exe, err := os.Executable()
	if err != nil {
		return resolved
	}

	exeDir := filepath.Dir(exe)
	candidates := []string{
		filepath.Join(exeDir, resolved),
		filepath.Join(filepath.Dir(exeDir), resolved),
		resolved,
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return candidates[0]
}

func findProc(name string) *syscall.Proc {
	proc, err := dll.FindProc(name)
	if err != nil {
		loadErr = fmt.Errorf("find proc %s: %w", name, err)
		return nil
	}
	return proc
}

func ensureDLL(dllPath string) error {
	loadOnce.Do(func() {
		var err error
		dll, err = syscall.LoadDLL(resolveDLLPath(dllPath))
		if err != nil {
			loadErr = fmt.Errorf("load dll: %w", err)
			return
		}

		procConnect = findProc("DM_Connect")
		procClose = findProc("DM_Close")
		procGetVersion = findProc("DM_GetVersion")
		procRestart = findProc("DM_RestartAppIfNecessary")
		procLastError = findProc("DM_GetLastError")
		procCheckForUpdates = findProc("DM_CheckForUpdates")
		procDownloadUpdate = findProc("DM_DownloadUpdate")
		procGetUpdateState = findProc("DM_GetUpdateState")
		procWaitUpdateState = findProc("DM_WaitForUpdateStateChange")
		procQuitAndInstall = findProc("DM_QuitAndInstall")
		procJsonToCanonical = findProc("DM_JsonToCanonical")
		procFreeString = findProc("DM_FreeString")

		procSetProductData = findProc("SetProductData")
		procSetProductId = findProc("SetProductId")
		procSetDataDirectory = findProc("SetDataDirectory")
		procSetDebugMode = findProc("SetDebugMode")
		procSetCustomDeviceFingerprint = findProc("SetCustomDeviceFingerprint")

		procSetLicenseKey = findProc("SetLicenseKey")
		procSetActivationMetadata = findProc("SetActivationMetadata")
		procActivateLicense = findProc("ActivateLicense")
		procActivateLicenseOffline = findProc("ActivateLicenseOffline")
		procGenerateOfflineDeactivationReq = findProc("GenerateOfflineDeactivationRequest")
		procGetLastActivationError = findProc("GetLastActivationError")

		procIsLicenseGenuine = findProc("IsLicenseGenuine")
		procIsLicenseValid = findProc("IsLicenseValid")
		procGetServerSyncGracePeriodExpiry = findProc("GetServerSyncGracePeriodExpiryDate")
		procGetActivationMode = findProc("GetActivationMode")

		procGetLicenseKey = findProc("GetLicenseKey")
		procGetLicenseExpiryDate = findProc("GetLicenseExpiryDate")
		procGetLicenseCreationDate = findProc("GetLicenseCreationDate")
		procGetLicenseActivationDate = findProc("GetLicenseActivationDate")
		procGetActivationCreationDate = findProc("GetActivationCreationDate")
		procGetActivationLastSyncedDate = findProc("GetActivationLastSyncedDate")
		procGetActivationId = findProc("GetActivationId")

		procGetLibraryVersion = findProc("GetLibraryVersion")
		procReset = findProc("Reset")
	})

	return loadErr
}
