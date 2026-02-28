package dmapi

import (
	gort "runtime"
	"syscall"
	"unsafe"
)

func callJSONEnvelopeWithOptions(proc *syscall.Proc, options map[string]interface{}) map[string]interface{} {
	request := marshalOptionalCJSON(options)
	var arg uintptr
	if len(request) > 0 {
		arg = uintptr(unsafe.Pointer(&request[0]))
	}

	ptr, _, _ := proc.Call(arg)
	gort.KeepAlive(request)
	if ptr == 0 {
		return nil
	}

	return parseEnvelope(ptrToOwnedString(ptr))
}

func (api *DmApi) CheckForUpdates(options map[string]interface{}) map[string]interface{} {
	return callJSONEnvelopeWithOptions(procCheckForUpdates, options)
}

func (api *DmApi) DownloadUpdate(options map[string]interface{}) map[string]interface{} {
	return callJSONEnvelopeWithOptions(procDownloadUpdate, options)
}

func (api *DmApi) CancelUpdateDownload(options map[string]interface{}) map[string]interface{} {
	return callJSONEnvelopeWithOptions(procCancelUpdateDownload, options)
}

func (api *DmApi) GetUpdateState() map[string]interface{} {
	ptr, _, _ := procGetUpdateState.Call()
	if ptr == 0 {
		return nil
	}
	return parseEnvelope(ptrToOwnedString(ptr))
}

func (api *DmApi) GetPostUpdateInfo() map[string]interface{} {
	ptr, _, _ := procGetPostUpdateInfo.Call()
	if ptr == 0 {
		return nil
	}
	return parseEnvelope(ptrToOwnedString(ptr))
}

func (api *DmApi) AckPostUpdateInfo(options map[string]interface{}) map[string]interface{} {
	return callJSONEnvelopeWithOptions(procAckPostUpdateInfo, options)
}

func (api *DmApi) WaitForUpdateStateChange(lastSequence uint64, timeoutMs uint32) map[string]interface{} {
	timeout := timeoutMs
	ptr, _, _ := procWaitUpdateState.Call(uintptr(lastSequence), uintptr(timeout))
	if ptr == 0 {
		return nil
	}
	return parseEnvelope(ptrToOwnedString(ptr))
}

func (api *DmApi) QuitAndInstall(options map[string]interface{}) int32 {
	request := marshalOptionalCJSON(options)
	var arg uintptr
	if len(request) > 0 {
		arg = uintptr(unsafe.Pointer(&request[0]))
	}

	ret, _, _ := procQuitAndInstall.Call(arg)
	gort.KeepAlive(request)
	return int32(ret)
}
