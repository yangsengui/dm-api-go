package dmapi

import (
	"encoding/json"
	"syscall"
	"unsafe"
)

func cStringBytes(s string) []byte {
	b := append([]byte(s), 0)
	if len(b) == 0 {
		return []byte{0}
	}
	return b
}

func cStringFromBuffer(buf []byte) string {
	for i, b := range buf {
		if b == 0 {
			return string(buf[:i])
		}
	}
	return string(buf)
}

func ptrToOwnedString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}

	var data []byte
	for i := 0; ; i++ {
		c := *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
		if c == 0 {
			break
		}
		data = append(data, c)
	}

	if procFreeString != nil {
		procFreeString.Call(ptr)
	}
	return string(data)
}

func ptrToStaticString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}

	var data []byte
	for i := 0; ; i++ {
		c := *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
		if c == 0 {
			break
		}
		data = append(data, c)
	}

	return string(data)
}

func callStatusBool(proc *syscall.Proc, args ...uintptr) bool {
	ret, _, _ := proc.Call(args...)
	return int32(ret) == 0
}

func callU32Out(proc *syscall.Proc) (uint32, bool) {
	var value uint32
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&value)))
	if int32(ret) != 0 {
		return 0, false
	}
	return value, true
}

func callStringOut(proc *syscall.Proc, size uint32) (string, bool) {
	if size == 0 {
		size = defaultBufferSize
	}

	buf := make([]byte, size)
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(size))
	if int32(ret) != 0 {
		return "", false
	}

	return cStringFromBuffer(buf), true
}

func normalizeOptions(options map[string]interface{}) map[string]interface{} {
	if options == nil {
		return map[string]interface{}{}
	}
	return options
}

func marshalCJSON(options map[string]interface{}) []byte {
	encoded, _ := json.Marshal(normalizeOptions(options))
	return append(encoded, 0)
}

func parseResponseData(raw string) map[string]interface{} {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return nil
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil
	}
	return data
}
