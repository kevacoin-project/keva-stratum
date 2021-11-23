package cnutil

// #cgo CFLAGS: -std=c11 -D_GNU_SOURCE
// #cgo LDFLAGS: -L${SRCDIR}/../build/cnutil -Wl,-rpath=\$ORIGIN/cnutil -lcnutil -Wl,-rpath ${SRCDIR} -lstdc++
// #include <stdlib.h>
// #include "src/cnutil.h"
import "C"
import (
	"encoding/hex"
	"log"
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

func Hash(blob []byte, fast bool, height int, seedHash string) []byte {
	output := make([]byte, 32)
	if fast {
		C.cryptonight_fast_hash((*C.char)(unsafe.Pointer(&blob[0])), (*C.char)(unsafe.Pointer(&output[0])), (C.uint32_t)(len(blob)))
	} else if len(seedHash) != 0 {
		// rx/keva
		seedHeight := C.rx_seedheight((C.uint64_t)(height))
		cnHash, err := hex.DecodeString(seedHash)
		if err != nil {
			log.Printf("Hash, failed to DecodeString: %v\n", err)
			return output
		}
		C.rx_slow_hash((C.uint64_t)(height), seedHeight, (*C.char)(unsafe.Pointer(&cnHash[0])),
			(*C.char)(unsafe.Pointer(&blob[0])), (C.uint64_t)(len(blob)), (*C.char)(unsafe.Pointer(&output[0])), 0, 0)
	} else {
		// cn/r
		C.cryptonight_hash((*C.char)(unsafe.Pointer(&blob[0])), (*C.char)(unsafe.Pointer(&output[0])), (C.uint32_t)(len(blob)), (C.int)(height))
	}
	return output
}

func FastHash(blob []byte) []byte {
	return Hash(blob, true, 0, "")
}
