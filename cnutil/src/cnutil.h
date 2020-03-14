#include <stdint.h>
#include <stddef.h>
#include "stdbool.h"

#ifdef __cplusplus
extern "C" {
#endif

uint32_t convert_blob(const char *blob, uint32_t len, char *out);
bool validate_address(const char *addr, uint32_t len);

void cryptonight_hash(const char* input, char* output, uint32_t len, int height);
void cryptonight_fast_hash(const char* input, char* output, uint32_t len);
uint64_t rx_seedheight(const uint64_t height);
void rx_slow_hash(const uint64_t mainheight, const uint64_t seedheight, const char *seedhash, const char *data, size_t length,
  char *hash, int miners, int is_alt);

#ifdef __cplusplus
}
#endif
