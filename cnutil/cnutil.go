package cnutil

// #cgo CFLAGS: -std=c11 -D_GNU_SOURCE
// #cgo LDFLAGS: -L${SRCDIR} -Wl,-rpath=\$ORIGIN/cnutil -lcnutil -Wl,-rpath ${SRCDIR} -lstdc++
// #include <stdlib.h>
// #include <string.h>
// #include "src/cnutil.h"
import "C"
import "unsafe"

func ConvertBlob(blob []byte) []byte {
	output := make([]byte, 76)
	out := (*C.char)(unsafe.Pointer(&output[0]))

	input := (*C.char)(unsafe.Pointer(&blob[0]))

	size := (C.uint32_t)(len(blob))
	C.convert_blob(input, size, out)
	return output
}

func GenerateAuxBlob(blob []byte) []byte {
	input := (*C.char)(unsafe.Pointer(&blob[0]))

	size := (C.uint32_t)(len(blob))
	var out *C.char
	blobSize := C.convert_blob_to_auxpow_blob(input, size, &out)
	defer C.free(unsafe.Pointer(out))
	output := make([]byte, blobSize)
	C.memcpy(unsafe.Pointer(&output[0]), unsafe.Pointer(out), (C.size_t)(blobSize))
	return output
}

func ValidateAddress(addr string) bool {
	input := C.CString(addr)
	defer C.free(unsafe.Pointer(input))

	size := (C.uint32_t)(len(addr))
	result := C.validate_address(input, size)
	return (bool)(result)
}

func Hash(blob []byte, fast bool) []byte {
	output := make([]byte, 32)
	if fast {
		C.cryptonight_fast_hash((*C.char)(unsafe.Pointer(&blob[0])), (*C.char)(unsafe.Pointer(&output[0])), (C.uint32_t)(len(blob)))
	} else {
		C.cryptonight_hash((*C.char)(unsafe.Pointer(&blob[0])), (*C.char)(unsafe.Pointer(&output[0])), (C.uint32_t)(len(blob)))
	}
	return output
}

func FastHash(blob []byte) []byte {
	return Hash(append([]byte{byte(len(blob))}, blob...), true)
}
