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
	defaultBufferSize   = 256
	defaultModeBufSize  = 64
	devLicenseErrorText = "Development license is missing or corrupted. Run `distromate sdk renew` to regenerate the dev certificate."

	envDmAPIPath          = "DM_API_PATH"
	envDmAppID            = "DM_APP_ID"
	envDmPublicKey        = "DM_PUBLIC_KEY"
	envDmLauncherEndpoint = "DM_LAUNCHER_ENDPOINT"
	envDmLauncherToken    = "DM_LAUNCHER_TOKEN"
)

const (
	dmErrOK uint32 = iota
	dmErrFail
	dmErrInvalidParameter
	dmErrAppIDNotSet
	dmErrLicenseKeyNotSet
	dmErrNotActivated
	dmErrLicenseExpired
	dmErrNetwork
	dmErrFileIO
	dmErrSignature
	dmErrBufferTooSmall
)

var activationErrorNames = map[uint32]string{
	dmErrOK:               "DM_ERR_OK",
	dmErrFail:             "DM_ERR_FAIL",
	dmErrInvalidParameter: "DM_ERR_INVALID_PARAMETER",
	dmErrAppIDNotSet:      "DM_ERR_APPID_NOT_SET",
	dmErrLicenseKeyNotSet: "DM_ERR_LICENSE_KEY_NOT_SET",
	dmErrNotActivated:     "DM_ERR_NOT_ACTIVATED",
	dmErrLicenseExpired:   "DM_ERR_LICENSE_EXPIRED",
	dmErrNetwork:          "DM_ERR_NETWORK",
	dmErrFileIO:           "DM_ERR_FILE_IO",
	dmErrSignature:        "DM_ERR_SIGNATURE",
	dmErrBufferTooSmall:   "DM_ERR_BUFFER_TOO_SMALL",
}

var (
	dll      *syscall.DLL
	loadErr  error
	loadOnce sync.Once

	procGetVersion                  *syscall.Proc
	procRestart                     *syscall.Proc
	procLastError                   *syscall.Proc
	procCheckForUpdates             *syscall.Proc
	procDownloadUpdate              *syscall.Proc
	procCancelUpdateDownload        *syscall.Proc
	procGetUpdateState              *syscall.Proc
	procGetPostUpdateInfo           *syscall.Proc
	procAckPostUpdateInfo           *syscall.Proc
	procWaitUpdateState             *syscall.Proc
	procQuitAndInstall              *syscall.Proc
	procJsonToCanonical             *syscall.Proc
	procFreeString                  *syscall.Proc
	procSetProductData              *syscall.Proc
	procSetProductId                *syscall.Proc
	procSetDataDirectory            *syscall.Proc
	procSetDebugMode                *syscall.Proc
	procSetCustomDeviceFP           *syscall.Proc
	procSetLicenseKey               *syscall.Proc
	procSetLicenseCallback          *syscall.Proc
	procActivateLicense             *syscall.Proc
	procGetLastActivationError      *syscall.Proc
	procIsLicenseGenuine            *syscall.Proc
	procIsLicenseValid              *syscall.Proc
	procGetServerSyncGrace          *syscall.Proc
	procGetActivationMode           *syscall.Proc
	procGetLicenseKey               *syscall.Proc
	procGetLicenseExpiryDate        *syscall.Proc
	procGetLicenseCreationDate      *syscall.Proc
	procGetLicenseActivationDate    *syscall.Proc
	procGetActivationCreationDate   *syscall.Proc
	procGetActivationLastSyncedDate *syscall.Proc
	procGetActivationID             *syscall.Proc
	procGetLibraryVersion           *syscall.Proc
	procReset                       *syscall.Proc
)

func resolveDLLPath(dllPath string) string {
	resolved := strings.TrimSpace(dllPath)
	if resolved == "" {
		resolved = strings.TrimSpace(os.Getenv(envDmAPIPath))
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
		if _, statErr := os.Stat(candidate); statErr == nil {
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

		procGetVersion = findProc("DM_GetVersion")
		procRestart = findProc("DM_RestartAppIfNecessary")
		procLastError = findProc("DM_GetLastError")
		procCheckForUpdates = findProc("DM_CheckForUpdates")
		procDownloadUpdate = findProc("DM_DownloadUpdate")
		procCancelUpdateDownload = findProc("DM_CancelUpdateDownload")
		procGetUpdateState = findProc("DM_GetUpdateState")
		procGetPostUpdateInfo = findProc("DM_GetPostUpdateInfo")
		procAckPostUpdateInfo = findProc("DM_AckPostUpdateInfo")
		procWaitUpdateState = findProc("DM_WaitForUpdateStateChange")
		procQuitAndInstall = findProc("DM_QuitAndInstall")
		procJsonToCanonical = findProc("DM_JsonToCanonical")
		procFreeString = findProc("DM_FreeString")

		procSetProductData = findProc("SetProductData")
		procSetProductId = findProc("SetProductId")
		procSetDataDirectory = findProc("SetDataDirectory")
		procSetDebugMode = findProc("SetDebugMode")
		procSetCustomDeviceFP = findProc("SetCustomDeviceFingerprint")

		procSetLicenseKey = findProc("SetLicenseKey")
		procSetLicenseCallback = findProc("SetLicenseCallback")
		procActivateLicense = findProc("ActivateLicense")
		procGetLastActivationError = findProc("GetLastActivationError")

		procIsLicenseGenuine = findProc("IsLicenseGenuine")
		procIsLicenseValid = findProc("IsLicenseValid")
		procGetServerSyncGrace = findProc("GetServerSyncGracePeriodExpiryDate")
		procGetActivationMode = findProc("GetActivationMode")

		procGetLicenseKey = findProc("GetLicenseKey")
		procGetLicenseExpiryDate = findProc("GetLicenseExpiryDate")
		procGetLicenseCreationDate = findProc("GetLicenseCreationDate")
		procGetLicenseActivationDate = findProc("GetLicenseActivationDate")
		procGetActivationCreationDate = findProc("GetActivationCreationDate")
		procGetActivationLastSyncedDate = findProc("GetActivationLastSyncedDate")
		procGetActivationID = findProc("GetActivationId")

		procGetLibraryVersion = findProc("GetLibraryVersion")
		procReset = findProc("Reset")
	})

	return loadErr
}
