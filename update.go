package dmapi

import "unsafe"

// Update methods (auto connect and close using DM_PIPE).

func (api *DmApi) CheckForUpdates(options map[string]interface{}) map[string]interface{} {
	request := marshalCJSON(options)
	return api.callPipeJSON(
		procCheckForUpdates,
		uintptr(unsafe.Pointer(&request[0])),
	)
}

func (api *DmApi) DownloadUpdate(options map[string]interface{}) map[string]interface{} {
	request := marshalCJSON(options)
	return api.callPipeJSON(
		procDownloadUpdate,
		uintptr(unsafe.Pointer(&request[0])),
	)
}

func (api *DmApi) GetUpdateState() map[string]interface{} {
	return api.callPipeJSON(procGetUpdateState)
}

func (api *DmApi) WaitForUpdateStateChange(lastSequence uint64, timeoutMs uint32) map[string]interface{} {
	return api.callPipeJSON(
		procWaitUpdateState,
		uintptr(lastSequence),
		uintptr(timeoutMs),
	)
}

func (api *DmApi) QuitAndInstall(options map[string]interface{}) bool {
	request := marshalCJSON(options)
	return api.callPipeAccepted(
		procQuitAndInstall,
		uintptr(unsafe.Pointer(&request[0])),
	)
}
