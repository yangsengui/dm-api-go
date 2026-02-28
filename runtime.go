package dmapi

import (
	"errors"
	"unsafe"
)

func (api *DmApi) GetVersion() string {
	ptr, _, _ := procGetVersion.Call()
	return ptrToStaticString(ptr)
}

func (api *DmApi) GetLastError() string {
	ptr, _, _ := procLastError.Call()
	return ptrToOwnedString(ptr)
}

func (api *DmApi) RestartAppIfNecessary() bool {
	ret, _, _ := procRestart.Call()
	return int32(ret) != 0
}

func JsonToCanonical(jsonStr string) (string, error) {
	if err := ensureDLL(""); err != nil {
		return "", err
	}

	request := cStringBytes(jsonStr)
	ptr, _, _ := procJsonToCanonical.Call(uintptr(unsafe.Pointer(&request[0])))
	if ptr == 0 {
		return "", errors.New("failed to convert to canonical format")
	}

	return ptrToOwnedString(ptr), nil
}

func (api *DmApi) JsonToCanonical(jsonStr string) (string, error) {
	return JsonToCanonical(jsonStr)
}
