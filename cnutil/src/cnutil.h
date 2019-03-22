#include <stdint.h>
#include "stdbool.h"

#ifdef __cplusplus
extern "C" {
#endif

uint32_t convert_blob(const char *blob, uint32_t len, char *out);

void cryptonight_hash(const char* input, char* output, uint32_t len, int height);
void cryptonight_fast_hash(const char* input, char* output, uint32_t len);

#ifdef __cplusplus
}
#endif
