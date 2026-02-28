# dm-api-go

Go SDK for DistroMate `dm_api` native library.

## Install

```bash
go get github.com/yangsengui/dm-api-go
```

## Quick Start (License)

```go
package main

import (
    "fmt"

    dmapi "github.com/yangsengui/dm-api-go"
)

func main() {
    api, err := dmapi.New("")
    if err != nil {
        panic(err)
    }

    api.SetProductData("<product-data>")
    api.SetProductId("your-product-id")
    api.SetLicenseKey("XXXX-XXXX-XXXX")

    if !api.ActivateLicense() {
        panic(api.GetLastError())
    }

    if !api.IsLicenseGenuine() {
        code, _ := api.GetLastActivationError()
        panic(fmt.Sprintf("license check failed: %s, err=%s", api.GetActivationErrorName(code), api.GetLastError()))
    }
}
```

## API Groups

- License setup: `SetProductData`, `SetProductId`, `SetDataDirectory`, `SetDebugMode`, `SetCustomDeviceFingerprint`
- License activation: `SetLicenseKey`, `SetLicenseCallback`, `ActivateLicense`, `GetLastActivationError`
- License state: `IsLicenseGenuine`, `IsLicenseValid`, `GetServerSyncGracePeriodExpiryDate`, `GetActivationMode`
- License details: `GetLicenseKey`, `GetLicenseExpiryDate`, `GetLicenseCreationDate`, `GetLicenseActivationDate`, `GetActivationCreationDate`, `GetActivationLastSyncedDate`, `GetActivationId`
- Update: `CheckForUpdates`, `DownloadUpdate`, `CancelUpdateDownload`, `GetUpdateState`, `GetPostUpdateInfo`, `AckPostUpdateInfo`, `WaitForUpdateStateChange`, `QuitAndInstall`
- General: `GetLibraryVersion`, `JsonToCanonical`, `GetLastError`, `Reset`

## Update API Notes

- Update APIs return parsed JSON envelope (`map[string]interface{}`) when transport succeeds.
- If native API returns `NULL`, Go SDK returns `nil`; check `GetLastError()`.
- `QuitAndInstall()` returns native `int32` status code directly:
  - `1`: accepted, process should exit soon
  - `-1`: business-level rejection (check `GetLastError()`)
  - `-2`: transport or parse error

## Environment Variables

- `DM_API_PATH`: optional path to native library
- `DM_APP_ID`, `DM_PUBLIC_KEY`: optional defaults for app identity
- `DM_LAUNCHER_ENDPOINT`, `DM_LAUNCHER_TOKEN`: launcher IPC variables used by update APIs

## Release

- CI runs `go test ./...` on Windows.
- Tag `v*` generates a release zip artifact.
- Go module publishing is tag-based (`go get` resolves by git tag).
