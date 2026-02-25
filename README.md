# dm-api-go

Go SDK for DistroMate `dm_api.dll`.

## Install

```bash
go get github.com/yangsengui/dm-api-go
```

## Integration Flow

1. Initialization: `SetProductData`, `SetProductId`.
2. Activation: `SetLicenseKey`, `ActivateLicense`.
3. Validation on startup: `IsLicenseGenuine` or `IsLicenseValid`.
4. Version/update: `GetVersion`, `GetLibraryVersion`, `CheckForUpdates`.

## Quick Start

```go
api, err := dmapi.New("")
if err != nil {
    panic(err)
}

api.SetProductData("<product_data>")
api.SetProductId("your-product-id", 0)
api.SetLicenseKey("XXXX-XXXX-XXXX")

if !api.ActivateLicense() {
    panic(api.GetLastError())
}
```

## Release

- CI runs `go test ./...` on Windows.
- Tag `v*` generates a release zip artifact.
- Go module publishing is tag-based (`go get` resolves by git tag).
