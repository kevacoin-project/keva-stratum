package cnutil

// #cgo CFLAGS: -std=c11 -D_GNU_SOURCE
// #cgo LDFLAGS: -L${SRCDIR} -Wl,-rpath=\$ORIGIN/cnutil -lcnutil -Wl,-rpath ${SRCDIR} -lstdc++
// #include <stdlib.h>
// #include "src/cnutil.h"
import "C"
import (
	"unsafe"
)

func ConvertBlob(blob []byte) []byte {
	output := make([]byte, 76)
	out := (*C.char)(unsafe.Pointer(&output[0]))

	input := (*C.char)(unsafe.Pointer(&blob[0]))

	size := (C.uint32_t)(len(blob))
	C.convert_blob(input, size, out)
	return output
}

func ValidateAddress(addr string) bool {
	input := C.CString(addr)
	defer C.free(unsafe.Pointer(input))

	size := (C.uint32_t)(len(addr))
	result := C.validate_address(input, size)
	return (bool)(result)
}

func Hash(blob []byte, fast bool, height int) []byte {
	output := make([]byte, 32)
	if fast {
		C.cryptonight_fast_hash((*C.char)(unsafe.Pointer(&blob[0])), (*C.char)(unsafe.Pointer(&output[0])), (C.uint32_t)(len(blob)))
	} else {
		C.cryptonight_hash((*C.char)(unsafe.Pointer(&blob[0])), (*C.char)(unsafe.Pointer(&output[0])), (C.uint32_t)(len(blob)), (C.int)(height))
	}
	return output
}

func FastHash(blob []byte) []byte {
	return Hash(append([]byte{byte(len(blob))}, blob...), true, 0)
}
