#include <stdint.h>
#include <string>
#include "cryptonote_basic/cryptonote_format_utils.h"
#include "common/base58.h"

using namespace cryptonote;

extern "C" uint32_t convert_blob(const char *blob, size_t len, char *out) {
    std::string input = std::string(blob, len);
    std::string output = "";

    block b = AUTO_VAL_INIT(b);
    if (!parse_and_validate_block_from_blob(input, b)) {
        return 0;
    }

    output = get_block_hashing_blob(b);
    output.copy(out, output.length(), 0);
    return output.length();
}
