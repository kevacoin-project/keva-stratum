cmake_minimum_required(VERSION 3.10)

project(assembler C ASM)
set(CXXLIB "cnutil")

find_package(OpenSSL REQUIRED)
find_package(Boost COMPONENTS thread system program_options date_time filesystem REQUIRED)

# Flags
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -D_GNU_SOURCE -maes")
set(CMAKE_C_FLAGS "-maes")

include_directories("src")
include_directories("src/contrib/epee/include")
include_directories("src/external/easylogging++")
include_directories("src/external/unbound")

# Build library
add_library(${CXXLIB} SHARED src/cnutil.cpp src/cryptonote_basic/cryptonote_format_utils.cpp
        src/crypto/tree-hash.c src/crypto/crypto.cpp src/crypto/crypto-ops.c src/crypto/crypto-ops-data.c
        src/crypto/hash.c src/crypto/slow-hash.c
        src/crypto/oaes_lib.c src/crypto/aesb.c
        src/crypto/hash-extra-groestl.c
        src/crypto/hash-extra-skein.c
        src/crypto/hash-extra-blake.c
        src/crypto/hash-extra-jh.c
        src/crypto/jh.c
        src/crypto/blake256.c
        src/crypto/random.c
        src/crypto/groestl.c
        src/crypto/skein.c
        src/crypto/chacha.c
        src/crypto/crypto_ops_builder/verify.c
        src/cryptonote_basic/account.cpp
        src/cryptonote_basic/cryptonote_basic_impl.cpp
        src/crypto/keccak.c src/common/base58.cpp
        src/crypto/CryptonightR_JIT.c
        src/contrib/epee/src/memwipe.c
        src/contrib/epee/src/mlocker.cpp
        src/contrib/epee/src/hex.cpp
        src/contrib/epee/src/wipeable_string.cpp
        src/ringct/rctTypes.cpp
        src/ringct/rctOps.cpp
        src/ringct/bulletproofs.cc
        src/ringct/multiexp.cc
        src/ringct/rctCryptoOps.c
        src/device/device.cpp
        src/device/device_default.cpp
        src/external/easylogging++/easylogging++.cc
        src/common/aligned.c
        src/common/perf_timer.cpp
        src/crypto/CryptonightR_template.S
        )

target_link_libraries(${CXXLIB}
    ${Boost_THREAD_LIBRARY}
    ${Boost_SYSTEM_LIBRARY}
    ${Boost_PROGRAM_OPTIONS_LIBRARY}
    ${Boost_DATE_TIME_LIBRARY}
    ${Boost_FILESYSTEM_LIBRARY}
    OpenSSL::SSL
)