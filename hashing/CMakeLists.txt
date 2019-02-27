set(LIB "hashing")

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -std=c11 -D_GNU_SOURCE")

add_library(${LIB} src/cryptonight.c src/cryptonight_fast.c src/crypto/aesb.c src/crypto/c_groestl.c src/crypto/c_keccak.c src/crypto/hash.c src/crypto/c_blake256.c src/crypto/c_jh.c src/crypto/c_skein.c src/crypto/oaes_lib.c)

