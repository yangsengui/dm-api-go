package dmapi

import (
	"runtime"
	"syscall"
	"unsafe"
)

func (api *DmApi) SetProductData(productData string) bool {
	value := cStringBytes(productData)
	return callStatusBool(procSetProductData, uintptr(unsafe.Pointer(&value[0])))
}

func (api *DmApi) SetProductId(productID string) bool {
	value := cStringBytes(productID)
	return callStatusBool(procSetProductId, uintptr(unsafe.Pointer(&value[0])))
}

func (api *DmApi) SetDataDirectory(directoryPath string) bool {
	value := cStringBytes(directoryPath)
	return callStatusBool(procSetDataDirectory, uintptr(unsafe.Pointer(&value[0])))
}

func (api *DmApi) SetDebugMode(enable bool) bool {
	var flag uint32
	if enable {
		flag = 1
	}
	return callStatusBool(procSetDebugMode, uintptr(flag))
}

func (api *DmApi) SetCustomDeviceFingerprint(fingerprint string) bool {
	value := cStringBytes(fingerprint)
	return callStatusBool(procSetCustomDeviceFP, uintptr(unsafe.Pointer(&value[0])))
}

func (api *DmApi) SetLicenseKey(licenseKey string) bool {
	value := cStringBytes(licenseKey)
	return callStatusBool(procSetLicenseKey, uintptr(unsafe.Pointer(&value[0])))
}

func (api *DmApi) SetLicenseCallback(callback func()) bool {
	if callback == nil {
		return false
	}

	thunk := syscall.NewCallback(func() uintptr {
		callback()
		return 0
	})

	ret, _, _ := procSetLicenseCallback.Call(thunk)
	if int32(ret) != 0 {
		return false
	}

	api.licenseCallbackFn = callback
	api.licenseCallbackAddr = thunk
	runtime.KeepAlive(api.licenseCallbackFn)
	return true
}

func (api *DmApi) ActivateLicense() bool {
	return callStatusBool(procActivateLicense)
}

func (api *DmApi) GetLastActivationError() (uint32, bool) {
	return callU32Out(procGetLastActivationError)
}

func (api *DmApi) IsLicenseGenuine() bool {
	return callStatusBool(procIsLicenseGenuine)
}

func (api *DmApi) IsLicenseValid() bool {
	return callStatusBool(procIsLicenseValid)
}

func (api *DmApi) GetServerSyncGracePeriodExpiryDate() (uint32, bool) {
	return callU32Out(procGetServerSyncGrace)
}

func (api *DmApi) GetActivationMode(bufferSize uint32) (string, string, bool) {
	if bufferSize == 0 {
		bufferSize = defaultModeBufSize
	}

	initial := make([]byte, bufferSize)
	current := make([]byte, bufferSize)
	ret, _, _ := procGetActivationMode.Call(
		uintptr(unsafe.Pointer(&initial[0])),
		uintptr(bufferSize),
		uintptr(unsafe.Pointer(&current[0])),
		uintptr(bufferSize),
	)
	if int32(ret) != 0 {
		return "", "", false
	}

	return cStringFromBuffer(initial), cStringFromBuffer(current), true
}

func (api *DmApi) GetLicenseKey(bufferSize uint32) (string, bool) {
	return callStringOut(procGetLicenseKey, bufferSize)
}

func (api *DmApi) GetLicenseExpiryDate() (uint32, bool) {
	return callU32Out(procGetLicenseExpiryDate)
}

func (api *DmApi) GetLicenseCreationDate() (uint32, bool) {
	return callU32Out(procGetLicenseCreationDate)
}

func (api *DmApi) GetLicenseActivationDate() (uint32, bool) {
	return callU32Out(procGetLicenseActivationDate)
}

func (api *DmApi) GetActivationCreationDate() (uint32, bool) {
	return callU32Out(procGetActivationCreationDate)
}

func (api *DmApi) GetActivationLastSyncedDate() (uint32, bool) {
	return callU32Out(procGetActivationLastSyncedDate)
}

func (api *DmApi) GetActivationId(bufferSize uint32) (string, bool) {
	return callStringOut(procGetActivationID, bufferSize)
}

func (api *DmApi) GetLibraryVersion() string {
	ptr, _, _ := procGetLibraryVersion.Call()
	return ptrToStaticString(ptr)
}

func (api *DmApi) Reset() bool {
	return callStatusBool(procReset)
}
