#include <cmath>
#include <stdint.h>
#include <string>
#include <algorithm>
#include "cryptonote_basic/cryptonote_basic.h"
#include "cryptonote_basic/cryptonote_format_utils.h"
#include "cryptonote_basic/blobdatatype.h"
#include "crypto/crypto.h"
#include "crypto/hash.h"
#include "common/base58.h"
#include "serialization/binary_utils.h"

using namespace cryptonote;

extern "C" uint32_t convert_blob(const char *blob, size_t len, char *out) {
    std::string input = std::string(blob, len);
    blobdata output = "";

    block b = AUTO_VAL_INIT(b);
    if (!cryptonote::parse_and_validate_block_from_blob(input, b)) {
        return 0;
    }

    output = cryptonote::get_block_hashing_blob(b);
    output.copy(out, output.length(), 0);
    return output.length();
}

extern "C" bool validate_address(const char *addr, size_t len) {
    std::string input = std::string(addr, len);
    std::string output = "";
    uint64_t prefix;
    return tools::base58::decode_addr(addr, prefix, output);
}

extern "C" void cryptonight_hash(const char* input, char* output, uint32_t len, int height) {
    const int variant = input[0] >= 7 ? input[0] - 6 : 0;
    crypto::cn_slow_hash(input, len, output, variant, 0, height);
}

extern "C" void cryptonight_fast_hash(const char* input, char* output, uint32_t len) {
    crypto::cn_fast_hash(input, len, output);
}